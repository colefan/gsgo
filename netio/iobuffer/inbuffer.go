package iobuffer

const (
	SIZE_1_K       int = 1024
	SIZE_10_K      int = 10240
	SIZE_20_K      int = 20480
	SIZE_64_K      int = 65536
	SIZE_1_M       int = 1048576
	MAX_WRITE_BUFF int = 65535
)

type InBuffer struct {
	Buff        []byte
	maxSize     int
	defaultSize int
}

func NewInBuffer(defaultCap, maxCap int) *InBuffer {
	return &InBuffer{defaultSize: defaultCap, maxSize: maxCap}
}

func (this *InBuffer) AppendData(bytes []byte) *InBuffer {
	l := len(bytes)
	if l > 0 {
		this.Buff = append(this.Buff, bytes...)
	}
	return this
}

func (this *InBuffer) GetBuffLen() int {
	return len(this.Buff)
}

func (this *InBuffer) CutPackData(packLen int) []byte {
	if packLen <= 0 {
		return nil
	}
	defer this.reduceCaps()
	//fmt.Println("cut pack data, packLen = ", packLen, "  raw buf = ", this.Buff)
	if len(this.Buff) >= packLen {
		tmp := this.Buff[0:packLen]
		//tmp = append(tmp, this.Buff[0:packLen]...)
		this.Buff = this.Buff[packLen:]
		return tmp
	} else {
		return nil
	}
}

func (this *InBuffer) reduceCaps() *InBuffer {
	nCap := cap(this.Buff)
	nLen := len(this.Buff)
	if nCap >= this.maxSize {
		if nCap/nLen >= 2 {
			//可以瘦身
			tmp := make([]byte, nLen, nLen)
			copy(tmp, this.Buff[0:nLen])
			this.Buff = tmp
		}

	}
	return this
}
