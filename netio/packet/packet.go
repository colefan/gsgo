/**
*packet def
*[LEN_16|CMDID_16|ID_32|FROM_16|TO_16|VCODE_16|PV_8|CSRC_8] = 16
*
*packet define file description
* @type = msg,entity
*<type = msg,cmd =110>
*<field =nid,type=uint32 desc=""/>
**/
package packet

import (
	"encoding/binary"
	"github.com/colefan/gsgo/netio/iobuffer"
)

const (
	PACKET_PROXY_HEADER_LEN = 16    //	代理协议的包头16个字节
	PACKET_LIMIT_SIZE       = 65535 //一个包的最大大小，包含包头，实际大小需要减掉16
)

//编码器
type Encoder interface {
	Encode(writeBuf *iobuffer.OutBuffer) *iobuffer.OutBuffer
}

//解码器
type Decoder interface {
	Decode(data []byte) (bool, []byte)
}

//实体解码器
type EntityDecoder interface {
	DecodeEntity(p *Packet) bool
}

type EntityEncoder interface {
	EncodeEntity(writeBuf *iobuffer.OutBuffer) *iobuffer.OutBuffer
}

//协议解码

//协议头
type Header struct {
	PackLen   uint16 //协议体的长度
	CmdID     uint16 //协议号ID
	ID        uint32 //用户ID
	FSID      uint16 //发送端服务ID
	TSID      uint16 //接受放服务ID
	ValidCode uint16 //校验码
	Version   uint8  //协议版本号 0-255
	ClientSrc uint8  //客户端来源
}

func (h *Header) Decode(data []byte) (bool, []byte) {
	if len(data) < PACKET_PROXY_HEADER_LEN {
		panic("not enough len for header,at least 16 bits")
		return false, nil
	}
	h.PackLen = binary.BigEndian.Uint16(data)
	data = data[2:]
	h.CmdID = binary.BigEndian.Uint16(data)
	data = data[2:]
	h.ID = binary.BigEndian.Uint32(data)
	data = data[4:]
	h.FSID = binary.BigEndian.Uint16(data)
	data = data[2:]
	h.TSID = binary.BigEndian.Uint16(data)
	data = data[2:]
	h.ValidCode = binary.BigEndian.Uint16(data)
	data = data[2:]
	h.Version = uint8(data[0])
	h.ClientSrc = uint8(data[1])
	data = data[2:]
	return false, data
}

func (h *Header) Encode(writeBuf *iobuffer.OutBuffer) *iobuffer.OutBuffer {
	if writeBuf == nil {
		writeBuf = iobuffer.NewOutBuffer(1024)
	}
	writeBuf.PutUint16(h.PackLen)
	writeBuf.PutUint16(h.CmdID)
	writeBuf.PutUint32(h.ID)
	writeBuf.PutUint16(h.FSID)
	writeBuf.PutUint16(h.TSID)
	writeBuf.PutUint16(h.ValidCode)
	writeBuf.PutUint8(h.Version)
	writeBuf.PutUint8(h.ClientSrc)
	return writeBuf
}

type Packet struct {
	Header
	RawData     []byte
	PackDecoded bool
}

func Packing(data []byte) *Packet {
	if len(data) < PACKET_PROXY_HEADER_LEN {
		return nil
	}
	pack := &Packet{RawData: data, PackDecoded: false}
	b := false
	b, pack.RawData = pack.Header.Decode(data)
	if !b {
		return nil
	} else {
		return pack
	}
}

func (this *Packet) IsDecoded() bool {
	return this.PackDecoded
}

//PACKET的解码方法，需要被子类重写
func (this *Packet) DecodePacket() bool {
	return false
}

//PACKET的编码方法需要被子类重写
func (this *Packet) EncodePacket(nLen int) *iobuffer.OutBuffer {
	return nil
}

func DecoderReadValue(this *Packet, v interface{}) bool {
	switch vtype := v.(type) {
	case *byte:
		*v.(*byte) = this.RawData[0]
		this.RawData = this.RawData[1:]
	case *uint16:
		*v.(*uint16) = binary.BigEndian.Uint16(this.RawData)
		this.RawData = this.RawData[2:]
	case *uint32:
		*v.(*uint32) = binary.BigEndian.Uint32(this.RawData)
		this.RawData = this.RawData[4:]
	case *uint64:
		*v.(*uint64) = binary.BigEndian.Uint64(this.RawData)
		this.RawData = this.RawData[8:]
	case *string:
		strLen := binary.BigEndian.Uint16(this.RawData)
		this.RawData = this.RawData[2:]
		if strLen > 0 && len(this.RawData) >= int(strLen) {
			*v.(*string) = string(this.RawData[0:strLen])
		} else {
			panic("not enough bytes to read for string")
			return false
		}
	case *[]byte:
		arrLen := binary.BigEndian.Uint16(this.RawData)
		this.RawData = this.RawData[2:]
		if arrLen > 0 {
			*v.(*[]byte) = append(*v.(*[]byte), this.RawData[0:arrLen]...)
		}
		this.RawData = this.RawData[arrLen:]
	case *[]uint16:
		arrLen := binary.BigEndian.Uint16(this.RawData)
		this.RawData = this.RawData[2:]
		if arrLen > 0 {
			for i := 0; i < int(arrLen); i++ {
				*v.(*[]uint16) = append(*v.(*[]uint16), binary.BigEndian.Uint16(this.RawData))
				this.RawData = this.RawData[2:]
			}
		}
	case *[]uint32:
		arrLen := binary.BigEndian.Uint16(this.RawData)
		this.RawData = this.RawData[2:]
		if arrLen > 0 {
			for i := 0; i < int(arrLen); i++ {
				*v.(*[]uint32) = append(*v.(*[]uint32), binary.BigEndian.Uint32(this.RawData))
				this.RawData = this.RawData[4:]
			}
		}
	case *[]uint64:
		arrLen := binary.BigEndian.Uint16(this.RawData)
		this.RawData = this.RawData[2:]
		if arrLen > 0 {
			for i := 0; i < int(arrLen); i++ {
				*v.(*[]uint64) = append(*v.(*[]uint64), binary.BigEndian.Uint64(this.RawData))
				this.RawData = this.RawData[8:]
			}
		}
	default:
		panic(vtype)
	}

	return true
}

func DecoderReadArrayLength(p *Packet) int {
	nLen := binary.BigEndian.Uint16(p.RawData)
	p.RawData = p.RawData[2:]
	return int(nLen)
}

//传entity实例时需要传入指针
func DecoderReadEntity(p *Packet, entity EntityDecoder) bool {
	return entity.DecodeEntity(p)
}
