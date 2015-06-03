package protocol_login2

import (
	"github.com/colefan/gsgo/netio/iobuffer"
	"github.com/colefan/gsgo/netio/packet"
)

const (
	CMD_C_LOGIN_REQUEST = 0x0001
)

type UserInfo struct {
	USER_ID	uint32 //
	Name	string //
	Pwd	string //
}

func (this *UserInfo) DecodeEntity(p *packet.Packet) bool {
	packet.DecoderReadValue(p, &this.USER_ID)
	packet.DecoderReadValue(p, &this.Name)
	packet.DecoderReadValue(p, &this.Pwd)
	return true
}

func (this *UserInfo) EncodeEntity(buf *iobuffer.OutBuffer) *iobuffer.OutBuffer {
	buf.PutRawValue(this.USER_ID)
	buf.PutRawValue(this.Name)
	buf.PutRawValue(this.Pwd)
	return buf
}
type LoginReq2 struct {
	*packet.Packet
	id	uint8 //
	account	string //
	account2	[]uint8 //
	UserData	UserInfo //
	UserLIst	[]UserInfo //
}

func (this *LoginReq2) DecodePacket() bool {
	if this.IsDecoded() {
		return true
	}
	packet.DecoderReadValue(this.Packet, &this.id)
	packet.DecoderReadValue(this.Packet, &this.account)
	packet.DecoderReadValue(this.Packet, &this.account2)
	packet.DecoderReadEntity(this.Packet, &this.UserData)
	arrLen:=packet.DecoderReadArrayLength(this.Packet)
	for i :=0; i < arrLen; i++ {
		e := &UserInfo{}
		packet.DecoderReadEntity(this.Packet, e)
		this.UserLIst = append(this.UserLIst, *e)
	}
	this.PackDecoded = true
	return true
}

func (this *LoginReq2) EncodePacket(nLen int) *iobuffer.OutBuffer {
	buf := iobuffer.NewOutBuffer(nLen)
	buf = this.Packet.Header.Encode(buf)
	buf.PutRawValue(this.id)
	buf.PutRawValue(this.account)
	buf.PutRawValue(this.account2)
	this.UserData.EncodeEntity(buf)
	if len(this.UserLIst) > 0 {
		buf.PutRawValue(uint16(len(this.UserLIst)))
		for _,tmp := range this.UserLIst {
			buf = tmp.EncodeEntity(buf)
		}
	}
	nPackLen := buf.GetLen() - packet.PACKET_PROXY_HEADER_LEN
	buf.SetUint16(uint16(nPackLen), 0)
	return buf
}
