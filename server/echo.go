package server

import (
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
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Error reading:", err)
	}
	log.Printf("Received: %s\n", buf[:n])
	conn.Write(buf[:n])
}
