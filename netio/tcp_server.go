package netio

import (
	"fmt"
	"github.com/colefan/gsgo/logs"
	"net"
	"strconv"
)

type TcpServerSocket struct {
	s        *Server
	address  string
	listener net.Listener
}

func NewTcpServerSocket(s *Server) *TcpServerSocket {
	server := &TcpServerSocket{s: s}
	server.address = ""
	return server

}

func (s *TcpServerSocket) Start() error {
	if len(s.address) == 0 {
		s.address = s.s.listenAdress + ":" + strconv.Itoa(int(s.s.listenPort))
	}
	var err error
	s.listener, err = net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("tcp listen address error " + s.address)
	}
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			logs.DefaultLogger.Error("listen accept error:", err)
			continue
		}
		go s.serve(conn)
	}

	return nil

}

func (s *TcpServerSocket) serve(c net.Conn) {
	session := NewConnection(c, s.s)
	session.Start()

}

func (s *TcpServerSocket) Close() error {
	return s.listener.Close()
}
