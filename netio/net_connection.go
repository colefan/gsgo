package netio

import (
	"github.com/colefan/gsgo/logs"
	"github.com/colefan/gsgo/netio/iobuffer"
	"io"
	"net"
	"sync"
)

type SessionHandlerInf interface {
	GetPackParser() PackParser
	GetPackDispatcher() PackDispatcher
}

//客户端连接接口
type ConnInf interface {
	Start()
	read()
	Write([]byte)
	write()
	Close()
	handshark()
}

type Connection struct {
	c              net.Conn          //物理连接
	s              SessionHandlerInf //所属服务器,或客户端
	handshaked     bool              //握手是否完成
	handshakecount int               //握手次数
	address        string
	readBuff       *iobuffer.InBuffer //读取缓存区
	writeBuff      *iobuffer.InBuffer //写入缓冲区
	writeMux       sync.Mutex
	status         int
}

func NewConnection(c net.Conn, s SessionHandlerInf) *Connection {
	return &Connection{c: c,
		s:         s,
		status:    SESSION_STATUS_INIT,
		readBuff:  iobuffer.NewInBuffer(iobuffer.SIZE_1_K, iobuffer.SIZE_64_K),
		writeBuff: iobuffer.NewInBuffer(iobuffer.SIZE_1_K, iobuffer.SIZE_1_M)}
}
func (this *Connection) Start() {
	//read first files

	this.address = this.c.RemoteAddr().String()
	this.status = SESSION_STATUS_OPEN
	//handshark
	if !this.handshaked {
		this.handshake()
	}
	//parsepack
	go this.read()
	go this.write()
}

func (this *Connection) read() {
	defer this.Close()
	for {
		var buf []byte
		n, err := this.c.Read(buf)

		if err != nil {
			logs.DefaultLogger.Error("Read data from conn err, err = ", err)
			break
		}

		data := this.s.GetPackParser().ParseMsg(buf, n, this.readBuff)
		if len(data) > 0 {
			go this.s.GetPackDispatcher().HandleMsg(data)
		}

	}
}

func (this *Connection) write() {
	defer this.Close()
	for {
		this.writeMux.Lock()
		if this.writeBuff.GetBuffLen() > 0 {
			tmp := this.writeBuff.CutPackData(this.writeBuff.GetBuffLen())
			if len(tmp) > 0 {
				_, err := this.c.Write(tmp)
				if err != nil {
					logs.DefaultLogger.Error("connection write error,", err)
					this.writeMux.Unlock()
					return
				}
			}
		}
		this.writeMux.Unlock()
	}
}

func (this *Connection) Write(data []byte) {
	this.writeMux.Lock()
	this.writeBuff.AppendData(data)
	this.writeMux.Unlock()
}

func (this *Connection) Close() {
	if this.status != SESSION_STATUS_CLOSED {
		err := this.c.Close()
		if err != nil {
			logs.DefaultLogger.Error("connection close error, ", err)
		}
		this.status = SESSION_STATUS_CLOSED
	}
}

func (this *Connection) handshake() {
	//过滤掉前面的128个字节,根据实际的需要进行修改
	filterBytes := make([]byte, 128)
	if !this.handshaked {
		n, err := io.ReadFull(this.c, filterBytes)
		if err != nil {
			logs.DefaultLogger.Error("handshake error")
		}
		println("n = ", n)
		this.handshaked = true
		this.handshakecount++
	}
}
