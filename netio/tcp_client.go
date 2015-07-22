package netio

import (
	"fmt"
	"net"
	"strconv"
)

type TcpClientSocket struct {
	c       *Client
	conn    *Connection
	address string
}

func NewTcpClientSocket(c *Client) *TcpClientSocket {
	myc := &TcpClientSocket{c: c}
	return myc

}

func (this *TcpClientSocket) Connect() error {
	if len(this.address) == 0 {
		this.address = this.c.serverAddress + ":" + strconv.Itoa(int(this.c.servePort))
	}
	fmt.Println("dial tcp ", this.address)
	conn, err := net.Dial("tcp", this.address)
	if err != nil {
		return err
	}

	session := NewConnection(conn, this.c)
	this.conn = session
	session.Start()
	this.c.reconnectchan <- false

	this.c.status = SESSION_STATUS_OPEN

	return nil
}

func (this *TcpClientSocket) Close() error {
	if this.c.status != SESSION_STATUS_CLOSED {
		this.conn.Close()
		this.c.status = SESSION_STATUS_CLOSED
		this.c.reconnectchan <- true
		this.c.reConnCount = 0
		this.c.StartReconnecting()
	}
	return nil
}

func (this *TcpClientSocket) Write(data []byte) {
	if this.conn != nil {
		this.conn.Write(data)
	} else {
		fmt.Println("client conn is empty")
	}
}
