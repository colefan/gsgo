package proxy

import (
	"github.com/colefan/gsgo/gameprotocol/protocol_proxy"
	"github.com/colefan/gsgo/netio"
	"github.com/colefan/gsgo/netio/packet"
)

type GsForwardRule struct {
}

func newGsForwardRule() *GsForwardRule {
	return &GsForwardRule{}
}

var gsForwarderInst = newGsForwardRule()

func (rule *GsForwardRule) FowardToServer(client *ProxyClient, cmdId uint16, pack *packet.Packet, conn netio.ConnInf) {
	server := client.GetGsRoute()
	if server == nil {
		errorResponse(cmdId, EC_GSROUTE_IS_NIL, pack.ID, conn)
		return
	}
	serverLink := client.GetPhysicalLink()
	if serverLink == nil {
		errorResponse(cmdId, EC_GSROUTE_IS_NIL, pack.ID, conn)
		return
	}
	if serverLink.GetBsStatus() != BS_STATUS_AUTHED {
		errorResponse(cmdId, EC_GSROUTE_IS_UNAUTHED, pack.ID, conn)
		return
	}
	serverLink.Write(pack.GetClientFromRawData())
}

func (rule *GsForwardRule) FowardToClient(cmdId uint16, pack *packet.Packet, conn netio.ConnInf) {
	clientID := pack.ID
	client := ClientManagerInst.FindClientByUserID(clientID)
	if client != nil {
		if cmdId == protocol_proxy.CMD_C_P_USER_OFFLINE_REQ {
			client.GetPhysicalLink().Close()
		} else {
			client.GetPhysicalLink().Write(pack.GetClientFromRawData())
		}
	} else {
		ProxyLog.Error("GS msg can't find physical net conn ,user_id = ", clientID)
	}

}
