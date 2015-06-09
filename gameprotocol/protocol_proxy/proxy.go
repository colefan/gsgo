package protocol_proxy

import (
	"github.com/colefan/gsgo/netio/iobuffer"
	"github.com/colefan/gsgo/netio/packet"
)

const (
	CMD_S_P_REG_REQ = 0x7F01
	CMD_S_P_REG_RESP = 0x7F02
	CMD_C_P_USER_OFFLINE_REQ = 0x7F03
	CMD_C_P_USER_OFFLINE_RESP = 0x7F04
	CMD_C_P_PROXY_ROUTE_REQ = 0x7F05
	CMD_C_P_PROXY_ROUTE_RESP = 0x7F06
	CMD_P_C_PROXY_ERROR_NT = 0x7F07
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
type ProxyClientOfflineReq struct {
	*packet.Packet
	UserID	uint32 //用户ID
}

func (this *ProxyClientOfflineReq) DecodePacket() bool {
	if this.IsDecoded() {
		return true
	}
	packet.DecoderReadValue(this.Packet, &this.UserID)
	this.PackDecoded = true
	return true
}

func (this *ProxyClientOfflineReq) EncodePacket(nLen int) *iobuffer.OutBuffer {
	buf := iobuffer.NewOutBuffer(nLen)
	buf = this.Packet.Header.Encode(buf)
	buf.PutRawValue(this.UserID)
	nPackLen := buf.GetLen() - packet.PACKET_PROXY_HEADER_LEN
	buf.SetUint16(uint16(nPackLen), 0)
	return buf
}
type ProxyClientOfflineResp struct {
	*packet.Packet
	UserID	uint32 //用户ID
}

func (this *ProxyClientOfflineResp) DecodePacket() bool {
	if this.IsDecoded() {
		return true
	}
	packet.DecoderReadValue(this.Packet, &this.UserID)
	this.PackDecoded = true
	return true
}

func (this *ProxyClientOfflineResp) EncodePacket(nLen int) *iobuffer.OutBuffer {
	buf := iobuffer.NewOutBuffer(nLen)
	buf = this.Packet.Header.Encode(buf)
	buf.PutRawValue(this.UserID)
	nPackLen := buf.GetLen() - packet.PACKET_PROXY_HEADER_LEN
	buf.SetUint16(uint16(nPackLen), 0)
	return buf
}
type ProxyRouteReq struct {
	*packet.Packet
}

func (this *ProxyRouteReq) DecodePacket() bool {
	if this.IsDecoded() {
		return true
	}
	this.PackDecoded = true
	return true
}

func (this *ProxyRouteReq) EncodePacket(nLen int) *iobuffer.OutBuffer {
	buf := iobuffer.NewOutBuffer(nLen)
	buf = this.Packet.Header.Encode(buf)
	nPackLen := buf.GetLen() - packet.PACKET_PROXY_HEADER_LEN
	buf.SetUint16(uint16(nPackLen), 0)
	return buf
}
type ProxyRouteResp struct {
	*packet.Packet
	Ip	string //可用IP
	Port	uint16 //可用PORT
	ExtStrVal	string //扩展数据
}

func (this *ProxyRouteResp) DecodePacket() bool {
	if this.IsDecoded() {
		return true
	}
	packet.DecoderReadValue(this.Packet, &this.Ip)
	packet.DecoderReadValue(this.Packet, &this.Port)
	packet.DecoderReadValue(this.Packet, &this.ExtStrVal)
	this.PackDecoded = true
	return true
}

func (this *ProxyRouteResp) EncodePacket(nLen int) *iobuffer.OutBuffer {
	buf := iobuffer.NewOutBuffer(nLen)
	buf = this.Packet.Header.Encode(buf)
	buf.PutRawValue(this.Ip)
	buf.PutRawValue(this.Port)
	buf.PutRawValue(this.ExtStrVal)
	nPackLen := buf.GetLen() - packet.PACKET_PROXY_HEADER_LEN
	buf.SetUint16(uint16(nPackLen), 0)
	return buf
}
type ProxyErrorNt struct {
	*packet.Packet
	ReqCmdID	uint16 //请求命令号
	ErrCode	uint16 //错误码
}

func (this *ProxyErrorNt) DecodePacket() bool {
	if this.IsDecoded() {
		return true
	}
	packet.DecoderReadValue(this.Packet, &this.ReqCmdID)
	packet.DecoderReadValue(this.Packet, &this.ErrCode)
	this.PackDecoded = true
	return true
}

func (this *ProxyErrorNt) EncodePacket(nLen int) *iobuffer.OutBuffer {
	buf := iobuffer.NewOutBuffer(nLen)
	buf = this.Packet.Header.Encode(buf)
	buf.PutRawValue(this.ReqCmdID)
	buf.PutRawValue(this.ErrCode)
	nPackLen := buf.GetLen() - packet.PACKET_PROXY_HEADER_LEN
	buf.SetUint16(uint16(nPackLen), 0)
	return buf
}
