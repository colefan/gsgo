package main

import (
	"fmt"
	"github.com/colefan/gsgo/console"
	"github.com/colefan/gsgo/gameprotocol/login"
	"github.com/colefan/gsgo/netio"
	"github.com/colefan/gsgo/netio/iobuffer"
	"github.com/colefan/gsgo/netio/packet"
	"runtime"
	"strconv"
	"time"
)


type MyClient struct {
	*netio.Client
}

func NewMyClient() *MyClient {
	c := &MyClient{}
	c.Client = netio.NewTcpClient()
	return c
}

func (this *MyClient) HandleMsg(cmdid uint16, pack *packet.Packet, conn netio.ConnInf) {
	//fmt.Println("[C]...read a msg, id = ", cmdid)
	buf := iobuffer.NewOutBuffer(int(pack.PackLen + packet.PACKET_PROXY_HEADER_LEN))
	buf = pack.Header.Encode(buf)
	for _, tmp := range pack.RawData {
		buf.PutByte(tmp)
	}

	nPackLen := buf.GetLen() - packet.PACKET_PROXY_HEADER_LEN
	buf.SetUint16(uint16(nPackLen), 0)

	//	if conn != nil {
	//		conn.Write(buf.GetData())
	//	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	//创建N个客户端
	clientList := make(map[int]*MyClient)
	for i := 0; i < 3000; i++ {
		go func() {
			c1 := NewMyClient()
			c1.SetPackDispatcher(netio.NewDefaultPackDispatcher())
			c1.SetPackParser(netio.NewDefaultParser())
			c1.SetServerAddress("192.168.15.26")
			c1.SetServerPort(12000)
			c1.GetPackDispatcher().AddPackEventListener(strconv.Itoa(i+1), c1)
			err := c1.Connect() //启动客户端，看乒乓能不能打起来，哈哈！
			if err != nil {
				fmt.Println("client error,", err)
			}
			clientList[i+1] = c1

			login := protocol_login.Login_Req{}
			login.Packet = packet.NewEmptyPacket()
			login.CmdID = 0x01
			login.UserName = "yjx"
			login.PWD = "1q2323"
			buf := login.EncodePacket(512)
			time.Sleep(20 * time.Millisecond)
			//fmt.Println("client send data :", buf.GetData())
			for {
				time.Sleep(time.Millisecond * 500)
				c1.Write(buf.GetData())
			}

		}()

	}

	console.CheckInput()

}
