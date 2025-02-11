package server

import (
	"io"
	"log"
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
	log.Println("Connected with...")
	log.Println(conn.RemoteAddr())
	io.Copy(conn, conn)
}
