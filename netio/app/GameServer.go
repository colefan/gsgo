package main

import (
	"fmt"
)

func main() {
	var k []byte
	fmt.Println(k)
	fmt.Println("Len = ", len(k))
	fmt.Println("Cap = ", cap(k))
	k = append(k, 1)
	k = append(k, 2)
	k = append(k, 3)
	var tmp []byte
	tmp = append(tmp, k[0:2]...)
	fmt.Println("tmp=>", tmp)
	fmt.Println("tmp len = ", len(tmp))
	fmt.Println("tmp cap = ", cap(tmp))
	fmt.Println("k=>", k)
	fmt.Println("Len 2 = ", len(k))
	fmt.Println("Cap 2 =", cap(k))
	k[0] = 10
	k[1] = 20
	k[2] = 30
	fmt.Println("k =>", k, "tmp=>", tmp)

	k = k[2:]
	fmt.Println(k)
	fmt.Println("Len 3 = ", len(k))
	fmt.Println("Cap 3 =", len(k))

	tmpList := make([]byte, 0)
	for i := 0; i < 129; i++ {
		tmpList = append(tmpList, byte(i))
	}

	fmt.Println("tmpList =>", tmpList, " len = ", len(tmpList), " cap = ", cap(tmpList))
	tmpList = tmpList[120:]
	fmt.Println("tmpList2 =>", tmpList, " len = ", len(tmpList), " cap = ", cap(tmpList))

	tmp2 := make([]byte, 3)
	copy(tmp2, tmpList[0:3])
	fmt.Println("tmp2 =>", tmp2, "len = ", len(tmp2), " cap = ", cap(tmp2))

}
