package tcpsocket

import (
	"github.com/colefan/gsgo/logs"
	"github.com/colefan/gsgo/netio"
	"net"
)

type TcpServerSocket struct {
	s       *netio.Server
	address string
}

func NewTcpServerSocket(s *netio.Server) *TcpServerSocket {
	server := &TcpServerSocket{s: s}
	server.address = ""
	return server

}

func (s *TcpServerSocket) Start() {
	listen, err := net.Listen("tcp", s.address)
	if err != nil {
		panic("tcp listen address error " + s.address)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			logs.DefaultLogger.Error("listen accept error:", err)
			continue
		}
		go s.serve(conn)
	}

}

func (s *TcpServerSocket) serve(c net.Conn) {
	session := netio.NewConnection(c, s.s)
	session.Start()

}
