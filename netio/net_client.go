package netio

import (
	"fmt"
	"sync"
	"time"
)

type ClientSocket interface {
	Connect() error
	Close() error
	Write(data []byte)
}
type Client struct {
	ClientSocket
	serverAddress   string
	servePort       uint16
	id              int
	status          int
	dispatcher      PackDispatcher
	parser          PackParser
	reConnCount     int
	bReconnable     bool
	shouldreconnect bool
	lockreconnect   sync.Mutex
}

func (this *Client) StartReconnecting() {

	if false == this.bReconnable {
		return
	}
	fmt.Println("reconnect...")
	this.lockreconnect.Lock()
	this.shouldreconnect = true
	this.reConnCount = 0
	this.lockreconnect.Unlock()

	go this.reconnecting()

}

func (this *Client) reconnecting() {

	for {
		this.reConnCount++
		fmt.Println("client is reconnecting....", this.reConnCount)
		err := this.Connect()
		if err != nil {
			time.Sleep(time.Second * 5)
		} else {
			fmt.Println("client is connected.")

			this.lockreconnect.Lock()
			this.shouldreconnect = false
			this.reConnCount = 0
			this.lockreconnect.Unlock()

			break
		}
	}

}

func (this *Client) Reconnectable() bool {
	return this.bReconnable
}

func (this *Client) SetReconnectable(b bool) {
	this.bReconnable = b
}

func (this *Client) SetServerAddress(address string) {
	this.serverAddress = address
}

func (this *Client) SetServerPort(port uint16) {
	this.servePort = port
}

func (this *Client) SetPackDispatcher(d PackDispatcher) {
	this.dispatcher = d
}

func (this *Client) GetPackDispatcher() PackDispatcher {
	return this.dispatcher
}

func (this *Client) SetPackParser(p PackParser) {
	this.parser = p
}

func (this *Client) GetPackParser() PackParser {
	return this.parser
}

func (this *Client) Connect() error {
	if len(this.serverAddress) == 0 && this.servePort == 0 {
		return fmt.Errorf("client connect error,server info invalid")
	}
	if this.ClientSocket == nil {
		return fmt.Errorf("clientsocket is nil")
	}
	return this.ClientSocket.Connect()
}

func (this *Client) Close() {
	if this.ClientSocket != nil {
		this.ClientSocket.Close()
	}
}

func (this *Client) Write(data []byte) {
	if this.ClientSocket != nil && len(data) > 0 {
		this.ClientSocket.Write(data)
	}

}

func NewTcpClient() *Client {
	c := &Client{}
	c.ClientSocket = NewTcpClientSocket(c)
	return c
}
