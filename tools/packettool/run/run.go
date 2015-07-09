package main

import (
	"fmt"
	"os"
	"strings"

	. "github.com/colefan/gsgo/tools/packettool"
	"github.com/colefan/gsgo/tools/utils"
)

//Useage:	packettool.exe xmlpath storepath
//
func main() {
	xmlPath := ""
	storePath := ""
	if len(os.Args) < 3 {
		fmt.Println("Useage:packettool.exe protocol_desc_file packet_store_path")
		xmlPath = "e:\\goproject\\src\\github.com\\colefan\\gsgo\\gameprotocol"
		storePath = "e:\\goproject\\src\\github.com\\colefan\\gsgo\\gameprotocol"
		fmt.Println("Packettool.exe will use default path : xmlpath = ", xmlPath, "   ,storepath = ", storePath)

	} else {
		xmlPath = strings.TrimSpace(os.Args[1])
		storePath = strings.TrimSpace(os.Args[2])
	}

	if len(xmlPath) < 0 {
		fmt.Println("protocol_desc_file is empty:(")
		return
	}

	if len(storePath) < 0 {
		fmt.Println("packet_store_path is empty:(")
		return
	}

	files, err := toolsutils.ListDir(xmlPath, ".xml")
	if err != nil {
		fmt.Println("xmlpath read error:", err)
		return
	}
	//fmt.Println("files,", len(files), "xmlpath=", xmlPath)
	for _, filepath := range files {
		//遍历目录，
		fmt.Println("filepath :", filepath)
		if perr := GenPacketForGoLang(storePath, filepath); perr != nil {
			fmt.Println("GenPacket error, filename =", filepath, ",error:", perr)
		}
	}
	return

}
