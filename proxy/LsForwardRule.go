package proxy

import (
	"github.com/colefan/gsgo/netio"
	"github.com/colefan/gsgo/netio/packet"
)

type LsForwardRule struct {
	ForwardRule
}

func newLsForwardRule() *LsForwardRule {
	return &LsForwardRule{}
}

var lsForwarderInst = newLsForwardRule()

func (rule *LsForwardRule) FowardToServer(client *ProxyClient, cmdID uint16, pack *packet.Packet, conn netio.ConnInf) {
	if code := rule.preRule(pack); code != 0 {
		errorResponse(cmdID, uint16(code), pack.ID, conn)
		return
	}
	server := client.GetLsRoute()
	if server == nil {
		errorResponse(cmdID, EC_LSROUTE_IS_NIL, pack.ID, conn)
		return
	}
	serverLink := server.GetPhysicalLink()
	if serverLink == nil {
		errorResponse(cmdID, EC_LSROUTE_IS_NIL, pack.ID, conn)
		return
	}
	serverLink.Write(pack.GetClientFromRawData())
	//rule.postRule(pack)//TODO 没有必要，但留有接口

}

func (rule *LsForwardRule) FowardToClient(client *ProxyClient, cmdID uint16, pack *packet.Packet, conn netio.ConnInf) {

}
