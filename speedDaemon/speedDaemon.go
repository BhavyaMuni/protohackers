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
	mu           sync.Mutex
	observations map[string][]Observation
	cameras      map[*net.Conn]Camera
	dispatchers  map[*net.Conn]Dispatcher
	tickets      map[uint16]chan *Ticket
	ticketDays   map[uint32]map[string]bool
}

func NewSpeedDaemonServer() *SpeedDaemonServer {
	ssd := &SpeedDaemonServer{
		observations: make(map[string][]Observation),
		cameras:      make(map[*net.Conn]Camera),
		dispatchers:  make(map[*net.Conn]Dispatcher),
		tickets:      make(map[uint16]chan *Ticket),
		ticketDays:   make(map[uint32]map[string]bool),
	}
	ssd.HandleConnectionFunc = ssd.handleConnection
	return ssd
}

func (ssd *SpeedDaemonServer) handleConnection(conn net.Conn) {
	buf := bufio.NewReader(conn)
	defer ssd.disconnectClient(conn)
	for {
		message, _, err := ParseMessage(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from client:", err)
			}
			break
		}
		message.Handle(ssd, &conn)
	}
}

func (ssd *SpeedDaemonServer) disconnectClient(conn net.Conn) {
	ssd.mu.Lock()
	defer ssd.mu.Unlock()

	delete(ssd.cameras, &conn)

	delete(ssd.dispatchers, &conn)

	log.Println("Disconnected client: ", conn.RemoteAddr())

	defer conn.Close()
}
