package netio

type PackParser interface {
	DecodeHead()
	DecodeBody()
}
