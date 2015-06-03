package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/colefan/gsgo/console"
	"github.com/colefan/gsgo/netio"
	"github.com/colefan/gsgo/netio/iobuffer"
	"github.com/colefan/gsgo/netio/packet"
	"github.com/colefan/gsgo/netio/qos"
)

type MyServer struct {
	*netio.Server
	rw        sync.Mutex
	nMsgCount int
}

func (this *MyServer) HandleMsg(cmdid uint16, pack *packet.Packet, conn netio.ConnInf) {
	//fmt.Println("[S]...read a msg, id = ", cmdid)
	buf := iobuffer.NewOutBuffer(int(pack.PackLen + packet.PACKET_PROXY_HEADER_LEN))
	buf = pack.Header.Encode(buf)
	for _, tmp := range pack.RawData {
		buf.PutByte(tmp)
	}

	nPackLen := buf.GetLen() - packet.PACKET_PROXY_HEADER_LEN
	buf.SetUint16(uint16(nPackLen), 0)

	if conn != nil {
		conn.Write(buf.GetData())
	}

}

func NewMyServer() *MyServer {
	s := &MyServer{}
	s.Server = netio.NewTcpSocketServer()
	return s
}

type ServerDispatcher struct {
	netio.DefaultPackDispatcher
	sessionMap map[int]int
	nums       int
}

func NewServerDispatcher() *ServerDispatcher {
	p := &ServerDispatcher{}
	p.sessionMap = make(map[int]int)
	return p
}

func (this *ServerDispatcher) SessionOpen(conn netio.ConnInf) {

	this.nums++
	fmt.Println("New clien coming:", this.nums)
}

func main() {
	fmt.Println("time = ", time.Now().Second())
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	s1 := NewMyServer()
	err := s1.Init(`{"ip":"127.0.0.1","port":12000}`)
	if err != nil {
		fmt.Println("net server init err,", err)
		return
	}

	s1.SetListenAddress("0.0.0.0")
	s1.SetListenPort(12000)

	s1.SetPackParser(netio.NewDefaultParser())

	s1.SetPackDispatcher(NewServerDispatcher())

	s1.GetPackDispatcher().AddPackEventListener("myserver", s1)
	s1.SetQos(netqos.NewServerQos())
	s1.GetQos().SetEnable(true)
	s1.GetQos().Stat()

	go s1.Start() //启动服务器

	console.CheckInput()
	s1.Close()
}
