package proxy

import (
	"github.com/colefan/gsgo/gameprotocol/protocol_proxy"
	"github.com/colefan/gsgo/netio"
	"github.com/colefan/gsgo/netio/packet"
)

//节点服务
//管理与各服务器之间的物理连接；
type NodeService struct {
	*netio.Server
	netio.DefaultPackDispatcher
}

func NewNodeService() *NodeService {
	s := &NodeService{}
	s.Server = netio.NewTcpSocketServer()
	return s
}

func (n *NodeService) InitService() error {
	n.SetListenAddress(ProxyConf.ListenIp)
	n.SetListenPort(uint16(ProxyConf.ListenPort))
	n.SetPackParser(netio.NewDefaultParser())
	n.SetPackDispatcher(n)
	n.GetPackDispatcher().AddPackEventListener("proxyserver", n)
	return nil
}

func (n *NodeService) SessionOpen(conn netio.ConnInf) {
	//TODO
	if conn == nil {
		panic("NodeService.SessionOpen(conn) error, conn is nil")
	}
	conn.SetBsStatus(BS_STATUS_OPENED)
	ProxyLog.Info("NodeService received a session: ", conn.GetRemoteIp())

}
func (n *NodeService) SessionClose(conn netio.ConnInf) {
	//TODO
	if conn.GetBsStatus() == BS_STATUS_AUTHED {
		NodeManager.UnRegNodeConnection(conn.GetConnID())
	}
	conn.SetBsStatus(BS_STATUS_CLOSED)
}

//收到来自各服务节点的消息
func (n *NodeService) HandleMsg(cmdID uint16, pack *packet.Packet, conn netio.ConnInf) {
	if conn == nil {
		panic("NodeService.HandleMsg error,conn is nil ")
	}
	switch conn.GetBsStatus() {
	case 0, BS_STATUS_OPENED:
		//需要先判断节点是已经通过了验证
		if protocol_proxy.CMD_S_S_REG_REQ == cmdID {
			node := NewServerNode()
			var regReq protocol_proxy.NodeRegReq
			regReq.Packet = pack
			if !regReq.DecodePacket() {
				ProxyLog.Error("invalid CMD_S_S_REG_REQ,Packet decode failed")
				conn.Close()
			} else {
				node.NodeType = regReq.NodeType
				node.GameId = regReq.GameId
				node.GameAreaId = regReq.GameAreaId
				node.GameServerId = regReq.GameServerId
				node.GameServerName = regReq.GameServerName
				node.GameServerDesc = regReq.GameServerDesc
				node.Ip = conn.GetRemoteIp()
				nRetCode := NodeManager.RegNodeConnection(node)
				if 0 == nRetCode {
					conn.SetBsStatus(BS_STATUS_AUTHED)
					conn.BindObj(node)
					//set connid;

				} else {
					ProxyLog.Info("server node register failed, ip = ", conn.GetRemoteIp(), " NodeType = ", node.NodeType, " IP = ", node.Ip, " errcode = ", nRetCode)
				}

				resp := protocol_proxy.NodeRegResp{}
				resp.Packet = packet.NewEmptyPacket()
				resp.Code = uint16(nRetCode)
				resp.CmdID = protocol_proxy.CMD_S_S_REG_RESP
				buf := resp.EncodePacket(256)
				conn.Write(buf.GetData())

			}

		} else {
			//不合法的请求，将其关闭
			ProxyLog.Error("invalid request ,cmdid = ", cmdID, " before server node is authed.")
			conn.Close()
		}
	case BS_STATUS_AUTHED:
	//请进入转发模式
	//fowardrule:::	TODO
	default:
		ProxyLog.Error("unknown server node status : ", conn.GetBsStatus())
	}

}
