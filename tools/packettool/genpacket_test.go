package packettool

import (
	"fmt"
	"testing"
)

func TestGenPacket(t *testing.T) {
	root, err := GenPacket("e:\\goproject\\src\\github.com\\colefan\\gsgo\\gameprotocol\\login.xml")
	if err != nil {
		fmt.Println("prase xml file error:", err)
	} else {
		root.GenGoProtocolFiles()
	}
}
