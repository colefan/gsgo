package protocol_proxy

import (
	"github.com/colefan/gsgo/netio/iobuffer"
	"github.com/colefan/gsgo/netio/packet"
)

const (
	CMD_S_S_REG_REQ = 0x7F01
	CMD_S_S_REG_RESP = 0x7F02
)

type NodeRegReq struct {
	*packet.Packet
	NodeType	uint16 //服务器节点类型：1-登录服；2-目录服；3-游戏逻辑服
	IP	string //
	Port	uint16 //
	GameId	uint32 //游戏ID,0为非游戏的服务器
	GameAreaId	uint32 //区域ID，支持分区分服
	GameServerId	uint32 //游戏服ID，同一游戏可以有多个游戏服，支持分区分服功能
	GameCode	string //游戏编码
	GameServerName	string //游戏服名称
	GameServerDesc	string //游戏服描述
}

func (this *NodeRegReq) DecodePacket() bool {
	if this.IsDecoded() {
		return true
	}
	packet.DecoderReadValue(this.Packet, &this.NodeType)
	packet.DecoderReadValue(this.Packet, &this.IP)
	packet.DecoderReadValue(this.Packet, &this.Port)
	packet.DecoderReadValue(this.Packet, &this.GameId)
	packet.DecoderReadValue(this.Packet, &this.GameAreaId)
	packet.DecoderReadValue(this.Packet, &this.GameServerId)
	packet.DecoderReadValue(this.Packet, &this.GameCode)
	packet.DecoderReadValue(this.Packet, &this.GameServerName)
	packet.DecoderReadValue(this.Packet, &this.GameServerDesc)
	this.PackDecoded = true
	return true
}

func (this *NodeRegReq) EncodePacket(nLen int) *iobuffer.OutBuffer {
	buf := iobuffer.NewOutBuffer(nLen)
	buf = this.Packet.Header.Encode(buf)
	buf.PutRawValue(this.NodeType)
	buf.PutRawValue(this.IP)
	buf.PutRawValue(this.Port)
	buf.PutRawValue(this.GameId)
	buf.PutRawValue(this.GameAreaId)
	buf.PutRawValue(this.GameServerId)
	buf.PutRawValue(this.GameCode)
	buf.PutRawValue(this.GameServerName)
	buf.PutRawValue(this.GameServerDesc)
	nPackLen := buf.GetLen() - packet.PACKET_PROXY_HEADER_LEN
	buf.SetUint16(uint16(nPackLen), 0)
	return buf
}
type NodeRegResp struct {
	*packet.Packet
	Code	uint16 //返回注册结果，默认为0
}

func (this *NodeRegResp) DecodePacket() bool {
	if this.IsDecoded() {
		return true
	}
	packet.DecoderReadValue(this.Packet, &this.Code)
	this.PackDecoded = true
	return true
}

func (this *NodeRegResp) EncodePacket(nLen int) *iobuffer.OutBuffer {
	buf := iobuffer.NewOutBuffer(nLen)
	buf = this.Packet.Header.Encode(buf)
	buf.PutRawValue(this.Code)
	nPackLen := buf.GetLen() - packet.PACKET_PROXY_HEADER_LEN
	buf.SetUint16(uint16(nPackLen), 0)
	return buf
}
