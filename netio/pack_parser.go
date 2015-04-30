package netio

import (
	"encoding/binary"
	"github.com/colefan/gsgo/netio/packet"
)

type PackParser interface {
	ParseMsg(readBuf []byte, n int, iobuff *IoBuffer) []byte
}

type DefaultParser struct {
}

func NewDefaultParser() *DefaultParser {
	return &DefaultParser{}
}

func (p *DefaultParser) ParseMsg(readBuf []byte, n int, iobuff *IoBuffer) []byte {
	if len(readBuf) >= n && n > 0 {
		iobuff.AppendData(readBuf[0:n])
		if len(iobuff.Buff) >= 2 {
			//头两个字节拿下来做为长度,默认支持代理模式
			packLen := binary.BigEndian.Uint16(iobuff.Buff)
			if iobuff.GetBuffLen() >= (int(packLen) + packet.PACKET_PROXY_HEADER_LEN) {
				data := iobuff.CutPackData(int(packLen) + packet.PACKET_PROXY_HEADER_LEN)
				return data
			}
		}
	}
	return nil
}
