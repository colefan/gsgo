package netio

import (
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
	conn, err := net.Dial("tcp", this.address)
	if err != nil {
		return err
	}

	session := NewConnection(conn, this.c)
	this.conn = session
	session.Start()
	this.c.status = SESSION_STATUS_OPEN

	return nil
}

func (this *TcpClientSocket) Close() error {
	if this.c.status != SESSION_STATUS_CLOSED {
		this.c.Close()
	}
	return nil
}
