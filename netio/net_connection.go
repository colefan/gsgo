package netio

import (
	"github.com/colefan/gsgo/logs"
	"io"
	"net"
)

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
	c              net.Conn //物理连接
	s              *Server  //所属服务器
	handshaked     bool     //握手是否完成
	handshakecount int      //握手次数
	address        string
	readBuf        []byte //读取缓存区
	writeBuf       []byte //写入缓冲区
}

func NewConnection(c net.Conn, s *Server) *Connection {
	return &Connection{c: c, s: s}
}
func (this *Connection) Start() {
	//read first files
	this.address = this.c.RemoteAddr().String()
	//handshark
	if !this.handshaked {
		this.handshake()
	}
	//parsepack
	go read()
	go write()
}

func (this *Connection) read() {
	defer this.Close()
	for {
		buf := make([]byte)
		n, err := this.c.Read(buf)
		if err != nil {
			logs.DefaultLogger.Error("Read data from conn err, err = ", err)
			break
		}

		data := this.s.GetPackParser().ParseMsg()
		if len(data) > 0 {
			go HandleMsg(data)
		}

	}
}

func (this *Connection) write() {
	defer this.Close()
}

func (this *Connection) Write(data []byte) {

}

func (this *Connection) Close() {

}

func (this *Connection) handshake() {
	//过滤掉前面的128个字节
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
