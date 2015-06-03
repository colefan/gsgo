package netio

import (
	"fmt"
	"net"
	"strconv"

	"github.com/colefan/gsgo/logs"
)

type TcpServerSocket struct {
	s           *Server
	address     string
	listener    net.Listener
	connIdIndex uint32
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
	defer func() {
		if s.listener != nil {
			s.listener.Close()
		}
	}()
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
		connId := s.GetNextConnID()
		go s.serve(connId, conn)
	}

	return nil

}

func (s *TcpServerSocket) GetNextConnID() uint32 {
	s.connIdIndex++
	if s.connIdIndex == 0 {
		s.connIdIndex++
	}
	return s.connIdIndex
}

func (s *TcpServerSocket) serve(id uint32, c net.Conn) {
	session := NewConnection(c, s.s)
	session.SetConnID(id)
	if s.s.GetQos() != nil {
		s.s.GetQos().StatAccpetConns()
		session.SetQos(s.s.GetQos())
	}
	s.s.GetPackDispatcher().SessionOpen(session)

	session.Start()

}

func (s *TcpServerSocket) Close() error {
	return s.listener.Close()
}
