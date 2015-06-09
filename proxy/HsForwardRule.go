package proxy

import (
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
	serverLink.Write(pack.GetClientFromRawData())

}

func (rule *HsForwardRule) FowardToClient(client *ProxyClient, cmdId uint16, pack *packet.Packet, conn netio.ConnInf) {

}
