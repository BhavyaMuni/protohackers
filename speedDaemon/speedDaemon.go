package speedDaemon

import (
	"bufio"
	"io"
	"log"
	"net"
	"sync"

	"github.com/BhavyaMuni/protohackers/server"
)

type SpeedDaemonServer struct {
	server.BaseServer
	mu              sync.Mutex
	observations    map[string][]Observation
	cameras         map[*net.Conn]Camera
	dispatchers     map[*net.Conn]Dispatcher
	roadDispatchers map[uint16]Dispatcher
}

func NewSpeedDaemonServer() *SpeedDaemonServer {
	ssd := &SpeedDaemonServer{
		observations:    make(map[string][]Observation),
		cameras:         make(map[*net.Conn]Camera),
		dispatchers:     make(map[*net.Conn]Dispatcher),
		roadDispatchers: make(map[uint16]Dispatcher),
	}
	ssd.HandleConnectionFunc = ssd.handleConnection
	return ssd
}

func (ssd *SpeedDaemonServer) handleConnection(conn net.Conn) {
	buf := bufio.NewReader(conn)
	for {
		message, messageType, err := ParseMessage(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from client:", err)
			}
			break
		}
		log.Println("Message:", message)
		log.Println("MessageType:", messageType)
		message.Handle(ssd, &conn)
	}
}
