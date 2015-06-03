//代理服务：
//主要功能：
//1)维护客户端的数据连接；
//2)转发工作服务器与客户端之间的通信协议；
//3)管理与具体工作服务器之间的物理连接;
//4)基本的网络安全策略;
//5)动态扩展：分布式节点的管理;
package proxy

import (
	"github.com/colefan/gsgo/netio"
	"github.com/colefan/gsgo/netio/packet"
)

//代理服务：
//1)收取来自客户端发送上来的消息；
//2)管理客户端上来的所有连接
type ProxyService struct {
	*netio.Server
	netio.DefaultPackDispatcher
}

func NewProxyService() *ProxyService {
	s := &ProxyService{}
	s.Server = netio.NewTcpSocketServer()
	return s
}

func (p *ProxyService) Start() error {
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

	if ok, node := NodeManagerInst.FindGsRoute(); ok {
		client.SetGsRoute(node)
	}

	ClientManagerInst.AddClient(client)
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
			//TODO

		} else {
			//TODO

		}
		conn.SetBsStatus(BS_STATUS_CLOSED)

	}

}

//收到消息的处理，需要知道是来自于client还是来自于服务器的消息
func (p *ProxyService) HandleMsg(cmdid uint16, pack *packet.Packet, conn netio.ConnInf) {
	//判断消息来源

}
