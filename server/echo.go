package server

import (
	"io"
	"net"
)

type EchoServer struct {
	BaseServer
}

func NewEchoServer() *EchoServer {
	es := &EchoServer{}
	es.HandleConnectionFunc = es.handleConnection
	return es
}

func (es EchoServer) handleConnection(conn net.Conn) {
	io.Copy(conn, conn)
}
