//代理服务：
//主要功能：
//1)维护客户端的数据连接；
//2)转发工作服务器与客户端之间的通信协议；
//3)管理与具体工作服务器之间的物理连接;
//4)基本的网络安全策略;
//5)动态扩展：分布式节点的管理;
package proxy

import (
	"sync"

	"github.com/colefan/gsgo/gameprotocol/protocol_proxy"
	"github.com/colefan/gsgo/netio"
	"github.com/colefan/gsgo/netio/packet"
)

//代理服务：
//1)收取来自客户端发送上来的消息；
//2)管理客户端上来的所有连接
type ProxyService struct {
	*netio.Server
	netio.DefaultPackDispatcher
	onlineClients int          //在线用户数
	onlineMutex   sync.RWMutex //用户读取锁
}

func NewProxyService() *ProxyService {
	s := &ProxyService{}
	s.Server = netio.NewTcpSocketServer()
	return s
}

func (p *ProxyService) InitService() error {
	return nil
}

//代理服务器接收一个物理链接
func (p *ProxyService) SessionOpen(conn netio.ConnInf) {
	if conn == nil {
		panic("ProxyService.SessionOpen(conn) error, conn is nil")
	}
	conn.SetBsStatus(BS_STATUS_OPENED)
	ProxyLog.Info("ProxyService received a session: ", conn.GetRemoteIp())

	client := NewProxyClient()
	client.SetPhysicalLink(conn)
	if ok, node := NodeManagerInst.FindLsRoute(); ok {
		client.SetLsRoute(node)
	}

	if ok, node := NodeManagerInst.FindHsRoute(); ok {
		client.SetHsRoute(node)
	}
	//	游戏服，需要等到真正去连接的时候才寻找路由
	//	if ok, node := NodeManagerInst.FindGsRoute(); ok {
	//		client.SetGsRoute(node)
	//	}

	ClientManagerInst.AddClient(client)
	p.IncreaseOnlines()
	ProxyLog.Debug("proxy session open , sessionid = ", conn.GetConnID(), " ip = ", conn.GetRemoteIp())

}

//代理服务发现一个物理链接被关闭
func (p *ProxyService) SessionClose(conn netio.ConnInf) {
	connId := conn.GetConnID()
	ProxyLog.Debug("proxy session closed , sessionid = ", connId, " ip = ", conn.GetRemoteIp())
	if connId == 0 {
		//没有进行验证，也没有记录到缓存中，可以直接断掉，不用做任何处理
		conn.SetBsStatus(BS_STATUS_CLOSED)
	} else {
		userId := conn.GetUID()
		if userId >= 0 {
			//用户已经进行验证，需要做如下处理
			//1.处理用户状态
			//2.通知其它服务器，该用户已经断线
			//3.维护用户列表
			if ok, client := ClientManagerInst.CloseClient(connId, userId); ok {
				gs := client.GetGsRoute()
				if gs != nil {
					//发送断线通知给游戏服
					req := &protocol_proxy.ProxyClientOfflineReq{}
					req.Packet = packet.NewEmptyPacket()
					req.CmdID = protocol_proxy.CMD_C_P_USER_OFFLINE_REQ
					req.UserID = userId
					req.FSID = NODE_TYPE_PROXY
					req.TSID = NODE_TYPE_GS
					req.ID = userId
					buf := req.EncodePacket(128)
					client.GetPhysicalLink().Write(buf.GetData())
				}

			}

		} else {
			ClientManagerInst.RemoveClientByConnID(connId)
		}

		conn.SetBsStatus(BS_STATUS_CLOSED)
	}
	p.DeduceOnlines()

}

//收到消息的处理，需要知道是来自于client还是来自于服务器的消息
func (p *ProxyService) HandleMsg(cmdid uint16, pack *packet.Packet, conn netio.ConnInf) {
	//判断消息来源
	nConnId := conn.GetConnID()
	nUserId := conn.GetUID()
	ProxyLog.Debug("from client sessionid = ", nConnId, ", user_id = ", nUserId, ", cmd_id = ", cmdid, " TSID = ", pack.TSID)
	var client *ProxyClient
	if nUserId > 0 {
		client = ClientManagerInst.FindClientByUserID(nUserId)
		if client == nil {
			ProxyLog.Error("from client msg can't find proxyclient instance, cmd_id = ", cmdid, ", user_id = ", nUserId, ", TSID = ", pack.TSID)
			conn.Close()
			return
		}

	} else {
		client = ClientManagerInst.FindClientByConnID(nConnId)
		if client == nil {
			ProxyLog.Error("from client msg can't find proxyclient instance, cmd_id = ", cmdid, ", conn_id = ", nConnId, ", TSID = ", pack.TSID)
			conn.Close()
			return
		}
	}

	if conn.GetBsStatus() != BS_STATUS_AUTHED {
		//没有认证，则需要判断，其是否过来了合理的消息
		if pack.TSID != NODE_TYPE_LS || pack.TSID != NODE_TYPE_PROXY {
			ProxyLog.Error("from client , before authed,cmd_id = ", cmdid, ", user_id = ", pack.ID, ", TSID = ", pack.TSID)
			conn.Close()
			return
		}

	}

	switch pack.TSID {
	case NODE_TYPE_LS:
		lsForwarderInst.FowardToServer(client, cmdid, pack, conn)
	case NODE_TYPE_HS:
		hsForwarderInst.FowardToServer(client, cmdid, pack, conn)
	case NODE_TYPE_GS:
		gsForwarderInst.FowardToServer(client, cmdid, pack, conn)
	case NODE_TYPE_PROXY:
		p.handleProxyLoadBalance(cmdid, pack, conn)
	default:
		ProxyLog.Error("from client unknow msg cmd_id = ", cmdid, ", user_id = ", nUserId, ", TSID = ", pack.TSID)
		conn.Close()

	}

}

//代理服务器本身的负载均衡功能，可以引导客户端连接到其它代理服务器
func (p *ProxyService) handleProxyLoadBalance(cmdid uint16, pack *packet.Packet, conn netio.ConnInf) {
	if protocol_proxy.CMD_C_P_PROXY_ROUTE_REQ != cmdid {
		ProxyLog.Error("from client, receive unknow msg fro proxyloadbanlance, cmd_id = ", cmdid, " , user_id = ", pack.ID)
		conn.Close()
		return
	}

	resp := &protocol_proxy.ProxyRouteResp{}
	resp.CmdID = protocol_proxy.CMD_C_P_PROXY_ROUTE_RESP
	resp.Packet = packet.NewEmptyPacket()

	if !ProxyConf.EnableReserveProxy {
		//不支持后备代理，则返回当前代理
		resp.Ip = ProxyConf.ListenIp
		resp.Port = uint16(ProxyConf.ListenPort)
	} else {
		if ProxyConf.LoadLimit > 0 && p.GetOnlines() >= ProxyConf.LoadLimit {
			//推荐尝试其它代理服务器
			resp.Ip = ProxyConf.ReserveProxyIp
			resp.Port = uint16(ProxyConf.ReserveProxyPort)
		} else {
			//推荐使用该代理服务器
			resp.Ip = ProxyConf.ListenIp
			resp.Port = uint16(ProxyConf.ListenPort)
		}

	}

	buf := resp.EncodePacket(128)
	conn.Write(buf.GetData())

}

func (p *ProxyService) GetOnlines() int {
	p.onlineMutex.RLock()
	defer p.onlineMutex.Unlock()
	return p.onlineClients
}

func (p *ProxyService) IncreaseOnlines() {
	p.onlineMutex.Lock()
	defer p.onlineMutex.Unlock()
	p.onlineClients++

}

func (p *ProxyService) DeduceOnlines() {
	p.onlineMutex.Lock()
	defer p.onlineMutex.Unlock()
	p.onlineClients--
	if p.onlineClients < 0 {
		p.onlineClients = 0
	}
}
