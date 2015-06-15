package protocol_comm

import (
	"github.com/colefan/gsgo/netio/iobuffer"
	"github.com/colefan/gsgo/netio/packet"
)

const (
	CMD_S_C_ERROR_NT = 0x00D2
	CMD_C_S_VCHECK_REQ = 0x00D3
	CMD_C_S_VCHECK_RESP = 0x00D4
)

type ServerErrorNt struct {
	*packet.Packet
	ReqCmdID	uint16 //请求命令号
	ErrCode	uint16 //错误码
}

func (this *ServerErrorNt) DecodePacket() bool {
	if this.IsDecoded() {
		return true
	}
	packet.DecoderReadValue(this.Packet, &this.ReqCmdID)
	packet.DecoderReadValue(this.Packet, &this.ErrCode)
	this.PackDecoded = true
	return true
}

func (this *ServerErrorNt) EncodePacket(nLen int) *iobuffer.OutBuffer {
	buf := iobuffer.NewOutBuffer(nLen)
	buf = this.Packet.Header.Encode(buf)
	buf.PutRawValue(this.ReqCmdID)
	buf.PutRawValue(this.ErrCode)
	nPackLen := buf.GetLen() - packet.PACKET_PROXY_HEADER_LEN
	buf.SetUint16(uint16(nPackLen), 0)
	return buf
}
type VersionCheckReq struct {
	*packet.Packet
	NodeType	uint16 //服务器类型
	GameID	uint32 //游戏ID
	GameCode	string //游戏编码
	Version	string //游戏版本号
}

func (this *VersionCheckReq) DecodePacket() bool {
	if this.IsDecoded() {
		return true
	}
	packet.DecoderReadValue(this.Packet, &this.NodeType)
	packet.DecoderReadValue(this.Packet, &this.GameID)
	packet.DecoderReadValue(this.Packet, &this.GameCode)
	packet.DecoderReadValue(this.Packet, &this.Version)
	this.PackDecoded = true
	return true
}

func (this *VersionCheckReq) EncodePacket(nLen int) *iobuffer.OutBuffer {
	buf := iobuffer.NewOutBuffer(nLen)
	buf = this.Packet.Header.Encode(buf)
	buf.PutRawValue(this.NodeType)
	buf.PutRawValue(this.GameID)
	buf.PutRawValue(this.GameCode)
	buf.PutRawValue(this.Version)
	nPackLen := buf.GetLen() - packet.PACKET_PROXY_HEADER_LEN
	buf.SetUint16(uint16(nPackLen), 0)
	return buf
}
type VersionCheckResp struct {
	*packet.Packet
	NodeType	uint16 //服务器类型
	GameID	uint32 //游戏ID
	GameCode	string //游戏编码
	Code	uint16 //错误码
}

func (this *VersionCheckResp) DecodePacket() bool {
	if this.IsDecoded() {
		return true
	}
	packet.DecoderReadValue(this.Packet, &this.NodeType)
	packet.DecoderReadValue(this.Packet, &this.GameID)
	packet.DecoderReadValue(this.Packet, &this.GameCode)
	packet.DecoderReadValue(this.Packet, &this.Code)
	this.PackDecoded = true
	return true
}

func (this *VersionCheckResp) EncodePacket(nLen int) *iobuffer.OutBuffer {
	buf := iobuffer.NewOutBuffer(nLen)
	buf = this.Packet.Header.Encode(buf)
	buf.PutRawValue(this.NodeType)
	buf.PutRawValue(this.GameID)
	buf.PutRawValue(this.GameCode)
	buf.PutRawValue(this.Code)
	nPackLen := buf.GetLen() - packet.PACKET_PROXY_HEADER_LEN
	buf.SetUint16(uint16(nPackLen), 0)
	return buf
}
