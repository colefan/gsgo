package protocol_login

import (
	"fmt"
	"testing"
)

func TestLoginMsg(t *testing.T) {
	loginReq := &Login_Req{}
	loginReq.CmdID = 0x01
	loginReq.UserName = "yjx"
	loginReq.PWD = "1234"
	loginReq.e.a = 10
	loginReq.e.b = 20
	loginReq.eList = append(loginReq.eList, loginReq.e)
	loginReq.eList = append(loginReq.eList, loginReq.e)

	buf := loginReq.EncodePacket(0)
	fmt.Println("buf =>", buf)
}
