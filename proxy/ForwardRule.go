package proxy

import (
	"github.com/colefan/gsgo/netio/packet"
)

type ForwardRule struct {
}

func (rule *ForwardRule) preRule(pack *packet.Packet) int {
	return 0
}

func (rule *ForwardRule) postRule(pack *packet.Packet) int {
	return 0
}
