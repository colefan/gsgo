package netio

import (
	"net"
)

type Connection struct {
	c              net.Conn //物理连接
	s              *Server  //所属服务器
	handshaked     bool     //握手是否完成
	handshakecount int      //握手次数
}

func NewConnection(c net.Conn, s *Server) *Connection {
	return &Connection{c: c, s: s}
}
func (this *Connection) Start() {

}

func (this *Connection) Read() {

}

func (this *Connection) Write() {

}
