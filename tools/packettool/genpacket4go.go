package packettool

import (
	"os"
	"strings"

	"github.com/colefan/gsgo/tools/utils"
)

func GenPacketForGoLang(storepath string, filename string) error {
	root, err := GenPacket(filename)
	if err != nil {
		return err
	}
	//根据root将文件存储到

	str, err := root.GenGoProtocolFiles()
	if err != nil {
		return err
	}
	if strings.LastIndex(storepath, string(os.PathSeparator)) == len(storepath)-1 {
		//fmt.Println("dddd")
		err := toolsutils.MakeFile(storepath+root.gopackage, toolsutils.ChFileExt(filename, ".go"), str)
		if err != nil {
			return err
		}
	} else {
		//fmt.Println("xxx:")
		err := toolsutils.MakeFile(storepath+string(os.PathSeparator)+root.gopackage, toolsutils.ChFileExt(toolsutils.GetFileName(filename), ".go"), str)
		if err != nil {
			return err
		}
	}

	return nil

}
