package packet

import (
	"fmt"
	"reflect"
	"testing"
)

type Login struct {
	username string
	userid   uint32
}

func TestPacket(t *testing.T) {
	pack := &Packet{}
	var login *Login
	login = &Login{}
	login.userid = 2
	pack.Body = login
	v := reflect.ValueOf(login)
	fmt.Println("v=>", v)
	for i := 0; i < v.Len(); i++ {
		fmt.Println(v.Field(i))

	}
	//fmt.Println(v)
	//println(pack.Body.userid)
}
