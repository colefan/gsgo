package main

import (
	"fmt"
	"github.com/colefan/gsgo/console"
	"github.com/colefan/gsgo/netio"
)

type ClientDispatcher struct {
	netio.DefaultPackDispatcher
}

func (this *ClientDispatcher) HandleMsg(data []byte) {
	nLen := len(data)
	if nLen <= 0 {
		return
	}

	pack := packet.Packing(data)

}

func main() {
	var s1 *netio.Server = netio.NewTcpSocketServer()
	err := s1.Init(`{"ip":"127.0.0.1","port":12000}`)
	if err != nil {
		fmt.Println("net server init err,", err)
		return
	}
	s1.SetListenAddress("0.0.0.0")
	s1.SetListenPort(12000)

	s1.SetPackParser(netio.NewDefaultParser())
	s1.SetPackDispatcher(netio.NewDefaultPackDispatcher())
	go s1.Start() //启动服务器

	c1 := netio.NewTcpClient()
	c1.SetPackDispatcher(netio.NewDefaultPackDispatcher())
	c1.SetPackParser(netio.NewDefaultParser())
	c1.SetServerAddress("127.0.0.1")
	c1.SetServerPort(12000)

	go c1.Connect() //启动客户端，看乒乓能不能打起来，哈哈！

	fmt.Println("other choice")

	console.CheckInput()

}
