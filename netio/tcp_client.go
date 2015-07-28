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
	this.c.status = SESSION_STATUS_OPEN
	this.c.GetPackDispatcher().SessionOpen(session)
	session.Start()
	this.c.lockreconnect.Lock()
	this.c.shouldreconnect = false
	this.c.lockreconnect.Unlock()

	return nil
}

func (this *TcpClientSocket) Close() error {
	if this.c.status != SESSION_STATUS_CLOSED {
		this.conn.Close()
		this.c.status = SESSION_STATUS_CLOSED

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
