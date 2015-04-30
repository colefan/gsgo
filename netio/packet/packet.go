/**
*packet def
*[LEN_16|CMDID_16|ID_32|FROM_16|TO_16|VCODE_16|PV_8|CSRC_8] = 16

*
**/
package packet

const (
	PACKET_PROXY_HEADER_LEN = 16    //	代理协议的包头16个字节
	PACKET_LIMIT_SIZE       = 65535 //一个包的最大大小，包含包头，实际大小需要减掉16
)

type header struct {
	PackLen   uint16 //协议体的长度
	CmdID     uint16 //协议号ID
	ID        uint32 //用户ID
	FSID      uint16 //发送端服务ID
	TSID      uint16 //接受放服务ID
	ValidCode uint16 //校验码
	Version   uint8  //协议版本号 0-255
	ClientSrc uint8  //客户端来源
}

type Packet struct {
	header
	data    []byte
	decoded bool
	Body    interface{}
}

func Pack(data []byte) *Packet {
	return nil

}
