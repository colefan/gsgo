package tcpsocket

import (
	"net"
)

type TcpServerSocket struct {
	listener *net.TCPListener
}
