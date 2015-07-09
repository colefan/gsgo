package packettool

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	TypePacket = 1
	TypeEntity = 2
)

const (
	O_BLANK_LINE   = "\n"
	O_IMPORT_LINES = "import (\n\t\"github.com/colefan/gsgo/netio/iobuffer\"\n\t\"github.com/colefan/gsgo/netio/packet\"\n)\n"
	O_SPACE        = " "
	O_TAB          = "\t"
)

type MyField struct {
	name      string
	fieldtype string
	value     string
	length    int
	comment   string
}

func NewMyField() *MyField {
	return &MyField{}
}

type MyPacket struct {
	classname string
	cmdname   string
	cmdvalue  string
	comment   string
	keys      []string
	fields    map[string]*MyField
	cctype    int
}

func NewMyPacket(cctype int) *MyPacket {
	return &MyPacket{cctype: cctype, fields: make(map[string]*MyField), keys: make([]string, 0)}
}

func (this *MyPacket) AddField(f *MyField) error {
	if len(f.name) <= 0 {
		return fmt.Errorf("field has no name")
	}

	tmp := this.fields[f.name]
	if tmp != nil {
		return fmt.Errorf("repeate field in same packet, pack-name=%s,field-name=%s", this.classname, f.name)
	}
	this.fields[f.name] = f
	this.keys = append(this.keys, f.name)
	return nil
}

func (this *MyPacket) GenGoPacketClass() string {
	// type A struct{
	//	*packet.Packet
	//	A	uint8
	//	B string
	//	C Entity
	//}
	str := "type " + this.classname + " struct {\n"
	str += O_TAB + "*packet.Packet\n"
	for _, key := range this.keys {
		f := this.fields[key]
		if f != nil {
			switch strings.ToLower(f.fieldtype) {
			case "uint8":
				str += O_TAB + f.name + O_TAB + "uint8 //" + f.comment + "\n"
			case "uint16":
				str += O_TAB + f.name + O_TAB + "uint16 //" + f.comment + "\n"
			case "uint32":
				str += O_TAB + f.name + O_TAB + "uint32 //" + f.comment + "\n"
			case "uint64":
				str += O_TAB + f.name + O_TAB + "uint64 //" + f.comment + "\n"
			case "string":
				str += O_TAB + f.name + O_TAB + "string //" + f.comment + "\n"
			case "entity":
				str += O_TAB + f.name + O_TAB + f.value + " //" + f.comment + "\n"
			case "entityarray":
				str += O_TAB + f.name + O_TAB + "[]" + f.value + " //" + f.comment + "\n"
			}

		}

	}
	str += "}\n"
	str += O_BLANK_LINE
	//decodepacke
	//func(this* A) DecodePacket() bool {
	//}
	str += "func (this *" + this.classname + ") DecodePacket() bool {\n"
	str += O_TAB + "if this.IsDecoded() {\n"
	str += O_TAB + O_TAB + "return true\n"
	str += O_TAB + "}\n"
	for _, key := range this.keys {
		f := this.fields[key]
		if f != nil {
			switch strings.ToLower(f.fieldtype) {
			case "entity":
				str += O_TAB + "packet.DecoderReadEntity(this.Packet, &this." + f.name + ")\n"
			case "entityarray":
				if f.value == "uint8" || f.value == "uint16" || f.value == "uint32" || f.value == "uint64" {
					str += O_TAB + "packet.DecoderReadValue(this.Packet, &this." + f.name + ")\n"
				} else {
					str += O_TAB + "arrLen:=packet.DecoderReadArrayLength(this.Packet)\n"
					str += O_TAB + "for i :=0; i < arrLen; i++ {\n"
					str += O_TAB + O_TAB + "e := &" + f.value + "{}\n"
					str += O_TAB + O_TAB + "packet.DecoderReadEntity(this.Packet, e)\n"
					str += O_TAB + O_TAB + "this." + f.name + " = append(this." + f.name + ", *e)\n"
					str += O_TAB + "}\n"

				}
			case "uint8", "uint16", "uint32", "uint64", "string":
				str += O_TAB + "packet.DecoderReadValue(this.Packet, &this." + f.name + ")\n"
			}

		}
	}
	str += O_TAB + "this.PackDecoded = true\n"
	str += O_TAB + "return true\n"
	str += "}\n"
	str += O_BLANK_LINE

	//func encodepacket
	str += "func (this *" + this.classname + ") EncodePacket(nLen int) *iobuffer.OutBuffer {\n"
	str += O_TAB + "buf := iobuffer.NewOutBuffer(nLen)\n"
	str += O_TAB + "buf = this.Packet.Header.Encode(buf)\n"
	for _, key := range this.keys {
		f := this.fields[key]
		if f != nil {
			switch strings.ToLower(f.fieldtype) {
			case "uint8", "uint16", "uint32", "uint64", "string":
				str += O_TAB + "buf.PutRawValue(this." + f.name + ")\n"
			case "entity":
				str += O_TAB + "this." + f.name + ".EncodeEntity(buf)\n"
			case "entityarray":
				if f.value == "uint8" || f.value == "uint16" || f.value == "uint32" || f.value == "uint64" {
					str += O_TAB + "buf.PutRawValue(this." + f.name + ")\n"
				} else {
					str += O_TAB + "if len(this." + f.name + ") > 0 {\n"
					str += O_TAB + O_TAB + "buf.PutRawValue(uint16(len(this." + f.name + ")))\n"
					str += O_TAB + O_TAB + "for _,tmp := range this." + f.name + " {\n"
					str += O_TAB + O_TAB + O_TAB + "buf = tmp.EncodeEntity(buf)\n"
					str += O_TAB + O_TAB + "}\n"
					str += O_TAB + "}\n"
				}

			}
		}
	}
	str += O_TAB + "nPackLen := buf.GetLen() - packet.PACKET_PROXY_HEADER_LEN\n"
	str += O_TAB + "buf.SetUint16(uint16(nPackLen), 0)\n"
	str += O_TAB + "return buf\n"
	str += "}"

	return str

}

func (this *MyPacket) GenGoEntityClass() string {
	str := "type " + this.classname + " struct {\n"
	for _, key := range this.keys {
		f := this.fields[key]
		if f != nil {
			switch strings.ToLower(f.fieldtype) {
			case "uint8":
				str += O_TAB + f.name + O_TAB + "uint8 //" + f.comment + "\n"
			case "uint16":
				str += O_TAB + f.name + O_TAB + "uint16 //" + f.comment + "\n"
			case "uint32":
				str += O_TAB + f.name + O_TAB + "uint32 //" + f.comment + "\n"
			case "uint64":
				str += O_TAB + f.name + O_TAB + "uint64 //" + f.comment + "\n"
			case "string":
				str += O_TAB + f.name + O_TAB + "string //" + f.comment + "\n"
			case "entity":
				str += O_TAB + f.name + O_TAB + f.value + " //" + f.comment + "\n"
			case "entityarray":
				str += O_TAB + f.name + O_TAB + "[]" + f.value + " //" + f.comment + "\n"
			}

		}

	}
	str += "}\n"
	str += O_BLANK_LINE
	//decodepacke
	//func(this* A) DecodePacket() bool {
	//}
	str += "func (this *" + this.classname + ") DecodeEntity(p *packet.Packet) bool {\n"
	for _, key := range this.keys {
		f := this.fields[key]
		if f != nil {
			switch strings.ToLower(f.fieldtype) {
			case "entity":
				str += O_TAB + "packet.DecoderReadEntity(p, &this." + f.name + ")\n"
			case "entityarray":
				if f.value == "uint8" || f.value == "uint16" || f.value == "uint32" || f.value == "uint64" {
					str += O_TAB + "packet.DecoderReadValue(p, &this." + f.name + ")\n"
				} else {
					str += O_TAB + "arrLen:=packet.DecoderReadArrayLength(p)\n"
					str += O_TAB + "for i :=0; i < arrLen; i++ {\n"
					str += O_TAB + O_TAB + "e := &" + f.value + "{}\n"
					str += O_TAB + O_TAB + "packet.DecoderReadEntity(p, e)\n"
					str += O_TAB + O_TAB + "this." + f.name + " = append(this." + f.name + ", *e)\n"
					str += O_TAB + "}\n"

				}
			case "uint8", "uint16", "uint32", "uint64", "string":
				str += O_TAB + "packet.DecoderReadValue(p, &this." + f.name + ")\n"
			}

		}
	}
	str += O_TAB + "return true\n"
	str += "}\n"
	str += O_BLANK_LINE

	//func encodepacket
	str += "func (this *" + this.classname + ") EncodeEntity(buf *iobuffer.OutBuffer) *iobuffer.OutBuffer {\n"

	for _, key := range this.keys {
		f := this.fields[key]
		if f != nil {
			switch strings.ToLower(f.fieldtype) {
			case "uint8", "uint16", "uint32", "uint64", "string":
				str += O_TAB + "buf.PutRawValue(this." + f.name + ")\n"
			case "entity":
				str += O_TAB + "this." + f.name + ".EncodeEntity(buf)\n"
			case "entityarray":
				if f.value == "uint8" || f.value == "uint16" || f.value == "uint32" || f.value == "uint64" {
					str += O_TAB + "buf.PutRawValue(this." + f.name + ")\n"
				} else {
					str += O_TAB + "if len(this." + f.name + ") > 0 {\n"
					str += O_TAB + O_TAB + "buf.PutRawValue(uint16(len(this." + f.name + ")))\n"
					str += O_TAB + O_TAB + "for _,tmp := range this." + f.name + " {\n"
					str += O_TAB + O_TAB + O_TAB + "buf = tmp.EncodeEntity(buf)\n"
					str += O_TAB + O_TAB + "}\n"
					str += O_TAB + "}\n"
				}

			}
		}
	}
	str += O_TAB + "return buf\n"
	str += "}"

	return str
}

type MyRoot struct {
	xmltype      string
	keys         []string
	packets      map[string]*MyPacket
	lastpackname string
	gopackage    string
}

func NewMyRoot() *MyRoot {
	return &MyRoot{packets: make(map[string]*MyPacket), keys: make([]string, 0)}
}

func (this *MyRoot) PutMyPacket(pack *MyPacket) bool {
	tmp := this.packets[pack.classname]
	if tmp != nil {
		return false
	}

	this.packets[pack.classname] = pack
	this.lastpackname = pack.classname
	this.keys = append(this.keys, pack.classname)
	return true
}

func (this *MyRoot) AddFieldInPack(field *MyField) error {
	if len(this.lastpackname) <= 0 {
		return fmt.Errorf("no packet in root")
	}

	tmp := this.packets[this.lastpackname]
	if tmp == nil {
		return fmt.Errorf("no pakcet in root")
	}

	return tmp.AddField(field)
}

func (this *MyRoot) GenGoProtocolFiles() (string, error) {
	stroutputH := ""
	stroutputH = "package " + this.gopackage + "\n"
	stroutputH += O_BLANK_LINE
	stroutputH += O_IMPORT_LINES
	stroutputConst := "const (\n"
	stroutputPacketContent := ""
	for _, key := range this.keys {
		item := this.packets[key]
		if item.cctype == TypeEntity {
			//stroutputConst += O_TAB + item.cmdname + O_SPACE + "=" + O_SPACE + item.cmdvalue + "\n"
			stroutputPacketContent += item.GenGoEntityClass() + "\n"

		} else if item.cctype == TypePacket {
			stroutputConst += O_TAB + item.cmdname + O_SPACE + "=" + O_SPACE + item.cmdvalue + "\n"
			stroutputPacketContent += item.GenGoPacketClass() + "\n"

		} else {
			return "", fmt.Errorf("unknow cctype :", item.cctype)
		}

	}

	stroutputConst += ")\n"
	//将文本内容保存的文件中去
	//TODO
	str := stroutputH + O_BLANK_LINE + stroutputConst + O_BLANK_LINE + stroutputPacketContent

	return str, nil
}

func GenPacket(filename string) (*MyRoot, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	decoder := xml.NewDecoder(bytes.NewBuffer(content))
	var t xml.Token
	var itemname string
	myroot := NewMyRoot()

	for t, err = decoder.Token(); err == nil; t, err = decoder.Token() {

		switch token := t.(type) {
		case xml.StartElement:
			itemname = token.Name.Local
			if "root" == strings.ToLower(itemname) {
				if len(myroot.xmltype) > 0 {
					return nil, fmt.Errorf("error: too many root flag in one xml file")
				}
				for _, attr := range token.Attr {
					//	if attr.Name.Local
					if "type" == strings.ToLower(attr.Name.Local) {
						myroot.xmltype = attr.Value
					} else if "gopackage" == strings.ToLower(attr.Name.Local) {
						myroot.gopackage = attr.Value
					}
				}

			} else if "pack" == strings.ToLower(itemname) {
				mypack := NewMyPacket(TypePacket)
				for _, attr := range token.Attr {
					switch strings.ToLower(attr.Name.Local) {
					case "classname":
						mypack.classname = attr.Value
					case "cmdname":
						mypack.cmdname = attr.Value
					case "cmdvalue":
						mypack.cmdvalue = attr.Value
					case "comment":
						mypack.comment = attr.Value
					default:
						fmt.Println("unhandle attibutes in pack:", attr.Name)
					}
				}

				if false == myroot.PutMyPacket(mypack) {

					return nil, fmt.Errorf("repeat pack in one xml file,pack name = %s", mypack.classname)
				}

			} else if "field" == strings.ToLower(itemname) {
				myfield := NewMyField()
				for _, attr := range token.Attr {
					switch strings.ToLower(attr.Name.Local) {
					case "name":
						myfield.name = attr.Value
					case "type":
						myfield.fieldtype = attr.Value
					case "value":
						myfield.value = attr.Value
					case "length":
						myfield.length, _ = strconv.Atoi(attr.Value)
					case "comment":
						myfield.comment = attr.Value
					default:
						fmt.Println("unhandle attibutes in field:", attr.Name)
					}

				}
				e1 := myroot.AddFieldInPack(myfield)
				if e1 != nil {
					return nil, e1
				}
			} else if "entity" == strings.ToLower(itemname) {
				mypack := NewMyPacket(TypeEntity)
				for _, attr := range token.Attr {
					switch strings.ToLower(attr.Name.Local) {
					case "classname":
						mypack.classname = attr.Value
					case "cmdname":
						mypack.cmdname = attr.Value
					case "cmdvalue":
						mypack.cmdvalue = attr.Value
					case "comment":
						mypack.comment = attr.Value
					default:
						fmt.Println("unhandle attibutes in pack:", attr.Name)
					}
				}

				if false == myroot.PutMyPacket(mypack) {
					return nil, fmt.Errorf("repeat pack in one xml file,pack name = %s", mypack.classname)
				}

			} else {
				fmt.Println("Unknow xml element in xml,please check:", itemname)
			}

		case xml.CharData:
			data := string([]byte(token))
			data = strings.TrimSpace(data)
			if len(data) > 0 {
				return nil, fmt.Errorf("Please use xml attributes not use xml value formate %s", data)
			}

		case xml.EndElement:

		}
	}

	return myroot, nil
}
