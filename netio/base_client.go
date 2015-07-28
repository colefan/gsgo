package netio

import (
	"fmt"
	"sync"

	"github.com/colefan/gsgo/netio/packet"
)

type BaseClient struct {
	*Client
	DefaultPackDispatcher
	bInited             bool
	socketStatus        int
	clientName          string
	serverListenAddress string
	serverListenPort    uint16
	mu                  sync.Mutex
}

func (this *BaseClient) InitClient(dispatcher PackDispatcher, listener PackListener, clientname string, serveraddress string, serverport uint16) error {
	this.clientName = clientname
	this.serverListenAddress = serveraddress
	this.serverListenPort = serverport
	if this.bInited {
		return nil
	}

	if this.Client == nil {
		this.Client = NewTcpClient()
	}
	this.SetServerAddress(this.serverListenAddress)
	this.SetServerPort(this.serverListenPort)
	this.SetPackDispatcher(dispatcher)
	this.SetPackParser(NewDefaultParser())
	this.AddPackEventListener("baseclientlistner", listener)
	if this.serverListenPort <= 0 {
		return fmt.Errorf("listen port invalid, port = ", this.serverListenPort)
	}
	return nil
}

func (this *BaseClient) GetSocketStatus() int {
	return this.socketStatus
}

func (this *BaseClient) SessionOpen(conn ConnInf) {
	fmt.Println("BaseClient:", this.clientName, ", session opened = ", conn.GetConnID(), " ip = ", conn.GetRemoteIp())
}
func (this *BaseClient) SessionClose(conn ConnInf) {
	fmt.Println("BaseClient:", this.clientName, ", session closed = ", conn.GetConnID())
}

func (this *BaseClient) HandleMsg(cmdID uint16, pack *packet.Packet, conn ConnInf) {
	//TODO 处理消息的地方，子类需要复写此方法
	fmt.Println("Received a msg ,cmd_id = ", cmdID)
}
