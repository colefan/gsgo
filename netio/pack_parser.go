package netio

type PackParser interface {
	ParseMsg(readBuf []byte, conn ConnInf) []byte
	DecodeMsg()
	DecodeHead()
	DecodeBody()
}
