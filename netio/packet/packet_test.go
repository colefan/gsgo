package packet

import (
	"encoding/binary"
	"fmt"
	"testing"
)

type Login struct {
	username string
	userid   uint32
}

type myinf interface {
	test()
}

type LoginReq struct {
	a     byte
	alist []byte
	b     Login
	blist []Login
}

func (this *LoginReq) test() {

}

func f1slice(b []int) {

	fmt.Println(b)

}

func putraw(data interface{}) {
	fmt.Println("be data=>", data)
	switch data := data.(type) {
	default:
		fmt.Println("af data =>", data)
	}
}

func findType(v interface{}) {
	switch vtype := v.(type) {
	default:
		if t, ok := v.([]myinf); true {
			//t[0].a = 11
			fmt.Println("t=>", t, "ok:", ok)
		}
		fmt.Println("vtype,", vtype)
	}
}

func TestPacket(t *testing.T) {

	if "BigEndian" == binary.BigEndian.String() {
		fmt.Println("big")
	} else {
		fmt.Println("small")
	}

	str1 := "我是姚"
	fmt.Println(len(str1))
	bs := []byte(str1)

	fmt.Println("Bs=>", bs, "len = ", len(bs))
	str2 := string(bs)
	fmt.Println("cs=>", str2)

	var b2 byte = 1
	var b3 uint8 = 2
	b3 = b2
	fmt.Println(b2, b3)

	login := &LoginReq{}
	fmt.Println("login.alist =>", login.alist)
	fmt.Println("login.blist=>", login.blist, &login.blist)

	var loginList []LoginReq
	fmt.Println("bbb=>", loginList)
	findType(loginList)
	list := make([]int32, 3)
	putraw(len(list))

	login

	//fmt.Println("loginReq=>", loginList)

}
