package protocol_login

import (
	"github.com/colefan/gsgo/netio/iobuffer"
	"github.com/colefan/gsgo/netio/packet"
)

const (
	CMD_C_LOGIN_REQ = 0x0001
	CMD_C_LOGIN_RESP = 0x0002
	CMD_C_LOGIN_VALID_REQ = 0x0003
	CMD_C_LOGIN_VALID_RESP = 0x0004
)

type LoginReq struct {
	*packet.Packet
	Account	string //账户
	CAID	uint8 //登录区域ID，预留的扩展字段
}

func (this *LoginReq) DecodePacket() bool {
	if this.IsDecoded() {
		return true
	}
	packet.DecoderReadValue(this.Packet, &this.Account)
	packet.DecoderReadValue(this.Packet, &this.CAID)
	this.PackDecoded = true
	return true
}

func (this *LoginReq) EncodePacket(nLen int) *iobuffer.OutBuffer {
	buf := iobuffer.NewOutBuffer(nLen)
	buf = this.Packet.Header.Encode(buf)
	buf.PutRawValue(this.Account)
	buf.PutRawValue(this.CAID)
	nPackLen := buf.GetLen() - packet.PACKET_PROXY_HEADER_LEN
	buf.SetUint16(uint16(nPackLen), 0)
	return buf
}
type LoginResp struct {
	*packet.Packet
	RandomCode	string //
}

func (this *LoginResp) DecodePacket() bool {
	if this.IsDecoded() {
		return true
	}
	packet.DecoderReadValue(this.Packet, &this.RandomCode)
	this.PackDecoded = true
	return true
}

func (this *LoginResp) EncodePacket(nLen int) *iobuffer.OutBuffer {
	buf := iobuffer.NewOutBuffer(nLen)
	buf = this.Packet.Header.Encode(buf)
	buf.PutRawValue(this.RandomCode)
	nPackLen := buf.GetLen() - packet.PACKET_PROXY_HEADER_LEN
	buf.SetUint16(uint16(nPackLen), 0)
	return buf
}
type LoginValidReq struct {
	*packet.Packet
	Account	string //帐号
	ValidCode	string //验证码
	CRandCode	string //随机码
}

func (this *LoginValidReq) DecodePacket() bool {
	if this.IsDecoded() {
		return true
	}
	packet.DecoderReadValue(this.Packet, &this.Account)
	packet.DecoderReadValue(this.Packet, &this.ValidCode)
	packet.DecoderReadValue(this.Packet, &this.CRandCode)
	this.PackDecoded = true
	return true
}

func (this *LoginValidReq) EncodePacket(nLen int) *iobuffer.OutBuffer {
	buf := iobuffer.NewOutBuffer(nLen)
	buf = this.Packet.Header.Encode(buf)
	buf.PutRawValue(this.Account)
	buf.PutRawValue(this.ValidCode)
	buf.PutRawValue(this.CRandCode)
	nPackLen := buf.GetLen() - packet.PACKET_PROXY_HEADER_LEN
	buf.SetUint16(uint16(nPackLen), 0)
	return buf
}
type LoginValidResp struct {
	*packet.Packet
	UFID	uint8 //用户来源ID
	UserId	uint32 //用户ID
	ValidFlag	string //登录凭证
	NeedRecover	uint8 //是否需要恢复现场
	GameAreaId	uint8 //游戏区域ID
	GameId	uint32 //游戏ID
	GameServerId	uint32 //游戏服ID
	GameServerIp	string //游戏服IP
	GameServerPort	uint16 //游戏服端口
}

func (this *LoginValidResp) DecodePacket() bool {
	if this.IsDecoded() {
		return true
	}
	packet.DecoderReadValue(this.Packet, &this.UFID)
	packet.DecoderReadValue(this.Packet, &this.UserId)
	packet.DecoderReadValue(this.Packet, &this.ValidFlag)
	packet.DecoderReadValue(this.Packet, &this.NeedRecover)
	packet.DecoderReadValue(this.Packet, &this.GameAreaId)
	packet.DecoderReadValue(this.Packet, &this.GameId)
	packet.DecoderReadValue(this.Packet, &this.GameServerId)
	packet.DecoderReadValue(this.Packet, &this.GameServerIp)
	packet.DecoderReadValue(this.Packet, &this.GameServerPort)
	this.PackDecoded = true
	return true
}

func (this *LoginValidResp) EncodePacket(nLen int) *iobuffer.OutBuffer {
	buf := iobuffer.NewOutBuffer(nLen)
	buf = this.Packet.Header.Encode(buf)
	buf.PutRawValue(this.UFID)
	buf.PutRawValue(this.UserId)
	buf.PutRawValue(this.ValidFlag)
	buf.PutRawValue(this.NeedRecover)
	buf.PutRawValue(this.GameAreaId)
	buf.PutRawValue(this.GameId)
	buf.PutRawValue(this.GameServerId)
	buf.PutRawValue(this.GameServerIp)
	buf.PutRawValue(this.GameServerPort)
	nPackLen := buf.GetLen() - packet.PACKET_PROXY_HEADER_LEN
	buf.SetUint16(uint16(nPackLen), 0)
	return buf
}
