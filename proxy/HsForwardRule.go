package proxy

import (
	"github.com/colefan/gsgo/gameprotocol/protocol_proxy"
	"github.com/colefan/gsgo/netio"
	"github.com/colefan/gsgo/netio/packet"
)

type HsForwardRule struct {
	ForwardRule
}

func newHsForwardRule() *HsForwardRule {
	return &HsForwardRule{}
}

var hsForwarderInst = newHsForwardRule()

func (rule *HsForwardRule) FowardToServer(client *ProxyClient, cmdId uint16, pack *packet.Packet, conn netio.ConnInf) {
	if code := rule.preRule(pack); code != 0 {
		errorResponse(cmdId, uint16(code), pack.ID, conn)
		return
	}
	server := client.GetHsRoute()
	if server == nil {
		errorResponse(cmdId, EC_HSROUTE_IS_NIL, pack.ID, conn)
		return
	}
	serverLink := server.GetPhysicalLink()
	if serverLink == nil {
		errorResponse(cmdId, EC_HSROUTE_IS_NIL, pack.ID, conn)
		return
	}
	if serverLink.GetBsStatus() != BS_STATUS_AUTHED {
		errorResponse(cmdId, EC_HSROUTE_IS_UNAUTHED, pack.ID, conn)
		return
	}
	serverLink.Write(pack.GetClientFromRawData())

}

func (rule *HsForwardRule) FowardToClient(cmdId uint16, pack *packet.Packet, conn netio.ConnInf) {
	clientID := pack.ID
	client := ClientManagerInst.FindClientByUserID(clientID)
	if client != nil {
		if cmdId == protocol_proxy.CMD_C_P_USER_OFFLINE_REQ {
			client.GetPhysicalLink().Close()
		} else {
			client.GetPhysicalLink().Write(pack.GetClientFromRawData())
		}

	} else {
		ProxyLog.Error("HS msg can't find phsical connection, user_id = ", clientID)
	}

}
