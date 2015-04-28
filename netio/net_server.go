//网络服务器接口
package netio

import (
	"encoding/json"
	"fmt"
	"net"
)

type ServerSocket interface {
	Start(address string)
	Close()
	serve(conn net.Conn)
	handshake(conn *Connection)
}

type Server struct {
	ServerSocket
	listenPort   uint16 `json:"port"`
	listenAdress string `json:"ip"`
	status       int
	dispatcher   PackDispatcher //消息分发器
	parser       PackParser     //消息解析器
}

func NewServer() *Server {
	return nil
}

func (s *Server) Init(config string) error {
	s.status = SERVER_STATUS_INIT
	if len(config) <= 0 {
		return fmt.Errorf("netio.Server.Init error,config is empty")
	}
	return json.Unmarshal([]byte(config), s)
}

func (s *Server) Start() error {

	//s.StartAndServe()
	return nil
}

func (s *Server) SetPackDispatcher(dispatcher PackDispatcher) {
	s.dispatcher = dispatcher
}

func (s *Server) GetPackDispatcher() PackDispatcher {
	return s.dispatcher
}

func (s *Server) SetPackParser(parser PackParser) {
	s.parser = parser
}

func (s *Server) GetPackParser() PackParser {
	return s.parser
}

func (s *Server) Shutdown() error {
	s.Close()
	return nil
}

func NewTcpSocketServer() *Server {
	return nil
}

func NewWebSocketServer() *Server {
	return nil
}
