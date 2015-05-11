package iobuffer

import (
	"encoding/binary"
)

type OutBuffer struct {
	data         []byte
	nLen         int
	littleEndian bool
}

func NewOutBuffer(size int) *OutBuffer {
	if size <= 0 {
		size = 1024
	}
	buf := &OutBuffer{data: make([]byte, size, size), nLen: 0, littleEndian: false}
	return buf
}

func (b *OutBuffer) GetData() []byte {
	return b.data[0:b.nLen]
}

func (b *OutBuffer) GetLen() int {
	return b.nLen
}
func (b *OutBuffer) SetEndian(endian string) {
	if endian == binary.LittleEndian.String() {
		b.littleEndian = true
	} else {
		b.littleEndian = false
	}
}

func (b *OutBuffer) PutByte(v byte) {
	if cap(b.data) > b.nLen {
		b.data[b.nLen] = v
		b.nLen++
	} else {
		b.data = append(b.data, v)
		b.nLen++
	}
}

func (b *OutBuffer) PutUint8(v uint8) {
	if cap(b.data) > b.nLen {
		b.data[b.nLen] = v
		b.nLen++
	} else {
		b.data = append(b.data, v)
		b.nLen++
	}
}

func (b *OutBuffer) PutUint16(v uint16) {
	tmp := make([]byte, 2)
	if b.littleEndian {
		binary.LittleEndian.PutUint16(tmp, v)
	} else {
		binary.BigEndian.PutUint16(tmp, v)
	}

	if cap(b.data) <= b.nLen+1 {
		b.data = append(b.data, 0, 0)
	}
	b.data[b.nLen] = tmp[0]
	b.nLen++
	b.data[b.nLen] = tmp[1]
	b.nLen++
}

func (b *OutBuffer) SetUint16(v uint16, nPos int) {
	tmp := make([]byte, 2)
	if b.littleEndian {
		binary.LittleEndian.PutUint16(tmp, v)
	} else {
		binary.BigEndian.PutUint16(tmp, v)
	}
	if b.nLen >= nPos+1 {
		b.data[nPos] = tmp[0]
		b.data[nPos+1] = tmp[1]
	}
}

func (b *OutBuffer) PutUint32(v uint32) {
	tmp := make([]byte, 4)
	if b.littleEndian {
		binary.LittleEndian.PutUint32(tmp, v)
	} else {
		binary.BigEndian.PutUint32(tmp, v)
	}

	if cap(b.data) <= b.nLen+3 {
		b.data = append(b.data, 0, 0, 0, 0)
	}
	for i := 0; i < 4; i++ {
		b.data[b.nLen+i] = tmp[i]
	}
	b.nLen += 4
}

func (b *OutBuffer) PutUint64(v uint64) {
	tmp := make([]byte, 8)
	if b.littleEndian {
		binary.LittleEndian.PutUint64(tmp, v)
	} else {
		binary.BigEndian.PutUint64(tmp, v)
	}

	if cap(b.data) <= b.nLen+7 {
		b.data = append(b.data, 0, 0, 0, 0, 0, 0, 0, 0)
	}

	for i := 0; i < 8; i++ {
		b.data[b.nLen+i] = tmp[i]
	}
	b.nLen += 8
}

func (b *OutBuffer) PutString(v string) {
	strLen := len(v)
	b.PutUint16(uint16(strLen))
	if strLen <= 0 {
		return
	}

	byteList := []byte(v)
	if cap(b.data) <= b.nLen+strLen {
		b.data = append(b.data, byteList...)
	}

	for i := 0; i < strLen; i++ {
		b.data[b.nLen+i] = byteList[i]
	}
	b.nLen += strLen
}

func (b *OutBuffer) PutRawValue(data interface{}) {
	switch data := data.(type) {
	case byte:
		b.PutByte(data)
	case uint16:
		b.PutUint16(data)
	case uint32:
		b.PutUint32(data)
	case uint64:
		b.PutUint64(data)
	case string:
		b.PutString(data)
	case []byte:
		//data = 1
		nLen := len(data)
		b.PutUint16(uint16(nLen))
		for i := 0; i < nLen; i++ {
			b.PutByte(data[i])
		}

	case []uint16:
		//data = 1
		nLen := len(data)
		b.PutUint16(uint16(nLen))
		for i := 0; i < nLen; i++ {
			b.PutUint16(data[i])
		}
	case []uint32:
		nLen := len(data)
		b.PutUint16(uint16(nLen))
		for i := 0; i < nLen; i++ {
			b.PutUint32(data[i])
		}
	case []uint64:
		nLen := len(data)
		b.PutUint16(uint16(nLen))
		for i := 0; i < nLen; i++ {
			b.PutUint64(data[i])
		}
	default:
		panic("undefine protocol field type")

	}
}
