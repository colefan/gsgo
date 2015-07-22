//网络服务器接口
package netio

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/colefan/gsgo/netio/qos"
)

type ServerSocket interface {
	Start() error
	Close() error
	serve(id uint32, conn net.Conn)
	//handshake(conn *Connection)

}

type Server struct {
	ServerSocket
	listenPort   uint16 `json:"port"`
	listenAdress string `json:"ip"`
	status       int
	dispatcher   PackDispatcher //消息分发器
	parser       PackParser     //消息解析器
	ClientNum    int
	qos          netqos.QosInf
}

func NewServer() *Server {
	return nil
}

func (s *Server) Init(config string) error {
	s.status = SERVER_STATUS_INITED
	if len(config) <= 0 {
		return fmt.Errorf("netio.Server.Init error,config is empty")
	}
	return json.Unmarshal([]byte(config), s)
}

func (s *Server) SetListenAddress(address string) {
	s.listenAdress = address
}

func (s *Server) SetListenPort(port uint16) {
	s.listenPort = port
}

func (s *Server) GetConfigJson() string {
	//return "`" + "{\"ip\":\"" + s.listenAdress + "\",\"port\":" + strconv.Itoa(int(s.listenPort)) + "}" + "`"
	return "{\"ip\":\"" + s.listenAdress + "\",\"port\":" + strconv.Itoa(int(s.listenPort)) + "}"
}
func (s *Server) Start() error {
	if s.status < SERVER_STATUS_INITED {
		panic("netio.server.start error, not inited.")
	}

	s.ServerSocket.Start()
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

func (s *Server) SetQos(qos netqos.QosInf) {
	s.qos = qos
}

func (s *Server) GetQos() netqos.QosInf {

	return s.qos
}

func NewTcpSocketServer() *Server {
	s := &Server{}
	s.ServerSocket = NewTcpServerSocket(s)
	return s
}

func NewWebSocketServer() *Server {
	return nil
}
