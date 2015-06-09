package proxy

import (
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
	serverLink.Write(pack.GetClientFromRawData())
}

func (rule *GsForwardRule) FowardToClient(client *ProxyClient, cmdId uint16, pack *packet.Packet, conn netio.ConnInf) {

}
