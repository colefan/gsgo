package protocol_login

import (
	"fmt"
	"testing"

	"github.com/colefan/gsgo/netio/packet"
)

func TestLoginMsg(t *testing.T) {

	loginReq := &Login_Req{}
	loginReq.Packet = packet.NewEmptyPacket()
	loginReq.CmdID = 0x01
	loginReq.UserName = "yjx"
	loginReq.PWD = "1234"
	loginReq.e.a = 10
	loginReq.e.b = 20
	loginReq.eList = append(loginReq.eList, loginReq.e)
	loginReq.eList = append(loginReq.eList, loginReq.e)

	buf := loginReq.EncodePacket(0)

	fmt.Println("login=>", loginReq)
	fmt.Println("buf =>", buf.GetData())
	pack := packet.Packing(buf.GetData())
	fmt.Println("cbuf=>", pack.GetClientFromRawData())

	loginReq2 := &Login_Req{}
	loginReq2.Packet = packet.Packing(buf.GetData())
	fmt.Println("packing=>", loginReq2.Packet)
	loginReq2.DecodePacket()
	fmt.Println("login2 =>", loginReq2)
}
