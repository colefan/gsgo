package protocol_login

import (
	"fmt"
	"github.com/colefan/gsgo/netio/iobuffer"
	"github.com/colefan/gsgo/netio/packet"
)

type Entity struct {
	a uint32
	b uint32
}

func (this *Entity) DecodeEntity(p *packet.Packet) bool {
	packet.DecoderReadValue(p, &this.a)
	packet.DecoderReadValue(p, &this.b)
	return true
}

func (this *Entity) EncodeEntity(writeBuf *iobuffer.OutBuffer) *iobuffer.OutBuffer {
	writeBuf.PutRawValue(this.a)
	writeBuf.PutRawValue(this.b)
	return writeBuf
}

type Login_Req struct {
	*packet.Packet
	UserName string //
	PWD      string //
	e        Entity
	eList    []Entity
}

type Login_Resp struct {
	PackData *packet.Packet
}

func (this *Login_Req) DecodePacket() bool {
	if this.IsDecoded() {
		return true
	}
	packet.DecoderReadValue(this.Packet, &this.UserName)
	fmt.Println("UserName = ", this.UserName)
	packet.DecoderReadValue(this.Packet, &this.PWD)
	fmt.Println("PWD =", this.PWD)
	packet.DecoderReadEntity(this.Packet, &this.e)

	arrLen := packet.DecoderReadArrayLength(this.Packet)
	for i := 0; i < arrLen; i++ {
		e := &Entity{}
		packet.DecoderReadEntity(this.Packet, e)
		this.eList = append(this.eList, *e)
	}
	this.PackDecoded = true
	return true
}

func (this *Login_Req) EncodePacket(nLen int) *iobuffer.OutBuffer {
	buf := iobuffer.NewOutBuffer(nLen)
	buf = this.Packet.Header.Encode(buf)
	buf.PutRawValue(this.UserName)
	buf.PutRawValue(this.PWD)
	buf = this.e.EncodeEntity(buf)
	if len(this.eList) > 0 {
		buf.PutRawValue(uint16(len(this.eList)))
		for _, tmp := range this.eList {
			buf = tmp.EncodeEntity(buf)
		}
	}
	nPackLen := buf.GetLen() - packet.PACKET_PROXY_HEADER_LEN
	buf.SetUint16(uint16(nPackLen), 0)

	return buf
}
