//网络服务器接口
package netio

import (
	"encoding/json"
	"fmt"
)

type NetServerSocket interface {
	Accept()
	Close()
	Start()
	Serve()
	StartAndServe()
	handshake()
}

type NetServer struct {
	*NetServerSocket
	listenPort   uint16 `json:"port"`
	listenAdress string `json:"ip"`
	status       int
}

func (s *NetServer) Init(config string) error {
	s.status = SERVER_STATUS_INIT
	if len(config) <= 0 {
		return fmt.Errorf("NetServer-Init-Error,config is empty")
	}
	return json.Unmarshal([]byte(config), s)
}

func (s *NetServer) Start() error {
	s.StartAndServe()
	return nil
}

func (s *NetServer) SetPackDispatcher(dispatcher PackDispatcher) {

}

func (s *NetServer) GetPackDispatcher() PackDispatcher {

}

func (s *NetServer) Shutdown() error {
	s.Close()
	return nil
}

func NewTcpSocketServer() *NetServer {
	return nil
}

func NewWebSocketServer() *NetServer {
	return nil
}
