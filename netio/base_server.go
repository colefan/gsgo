package netio

import (
	"fmt"
	"sync"

	"github.com/colefan/gsgo/netio/packet"
)

//提供完整的最基础的服务器
type BaseServer struct {
	*Server
	DefaultPackDispatcher
	bInited    bool
	serverName string
	mu         sync.RWMutex //读写锁
}

//初始化服务器
func (this *BaseServer) InitServer(dispatcher PackDispatcher, listener PackListener, servername string, ip string, port uint16) (err error) {

	if this.bInited {
		return nil
	}
	this.serverName = servername
	if this.Server == nil {
		this.Server = NewTcpSocketServer()
	}
	this.SetListenAddress(ip)
	this.SetListenPort(port)
	this.SetPackParser(NewDefaultParser())
	this.SetPackDispatcher(dispatcher)
	this.AddPackEventListener("baseserverlistener", listener)
	err = this.Init(this.GetConfigJson())
	if err == nil {
		this.bInited = true
	}
	return err
}

//运行服务器
func (this *BaseServer) Run() {
	go this.Server.Start()
}

//关闭服务器
func (this *BaseServer) ShutDown() error {
	return this.Server.Shutdown()
}

func (this *BaseServer) SessionOpen(conn ConnInf) {
	fmt.Println("BaseServer:", this.serverName, ", opened a session id = ", conn.GetConnID(), " ip = ", conn.GetRemoteIp())
}
func (this *BaseServer) SessionClose(conn ConnInf) {
	fmt.Println("BaseServer:", this.serverName, ", closed a session id = ", conn.GetConnID())
}

func (this *BaseServer) HandleMsg(cmdID uint16, pack *packet.Packet, conn ConnInf) {
	//TODO 处理消息的地方，子类需要复写此方法
	fmt.Println("Received a msg ,cmd_id = ", cmdID)
}
