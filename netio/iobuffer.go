package netio

const (
	SIZE_1_K  int = 1024
	SIZE_10_K int = 10240
	SIZE_64_K int = 65536
	SIZE_1_M  int = 1048576
)

type IoBuffer struct {
	buff        []byte
	maxSize     int
	defaultSize int
	startPos int
	endPos int
}

func NewIoBuffer(defaultCap, maxCap int) *IoBuffer {
	return &IoBuffer{defaultSize: defaultCap, maxSize: maxCap}
}

func
