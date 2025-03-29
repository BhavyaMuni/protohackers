package speedDaemon

import (
	"bufio"
	"bytes"
	"encoding/binary"
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
	heartbeats   map[*net.Conn]bool
	ticketDays   map[uint32]map[string]bool
}

func NewSpeedDaemonServer() *SpeedDaemonServer {
	ssd := &SpeedDaemonServer{
		observations: make(map[string][]Observation),
		cameras:      make(map[*net.Conn]Camera),
		dispatchers:  make(map[*net.Conn]Dispatcher),
		tickets:      make(map[uint16]chan *Ticket),
		heartbeats:   make(map[*net.Conn]bool),
		ticketDays:   make(map[uint32]map[string]bool),
	}
	ssd.HandleConnectionFunc = ssd.handleConnection
	return ssd
}

func (ssd *SpeedDaemonServer) handleConnection(conn net.Conn) {
	buf := bufio.NewReader(conn)
	defer ssd.disconnectClient(&conn)
	for {
		message, _, err := ParseMessage(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from client: ", err)
				ssd.SendError(&conn, "Error reading from client")
			}
			break
		}
		log.Println("Received message: ", message, "from", conn.RemoteAddr(), "conn", &conn)
		message.Handle(ssd, &conn)
	}
}

func (ssd *SpeedDaemonServer) disconnectClient(conn *net.Conn) {
	ssd.mu.Lock()
	defer ssd.mu.Unlock()
	defer (*conn).Close()

	delete(ssd.cameras, conn)

	delete(ssd.dispatchers, conn)

	delete(ssd.heartbeats, conn)

	log.Println("Disconnected client: ", (*conn).RemoteAddr())

}

func (ssd *SpeedDaemonServer) SendError(conn *net.Conn, message string) {
	defer ssd.disconnectClient(conn)
	messageType := ErrorMessageType
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, messageType)
	binary.Write(buf, binary.BigEndian, uint8(len(message)))
	buf.WriteString(message)
	err := binary.Write(*conn, binary.BigEndian, buf.Bytes())
	if err != nil {
		log.Println("Error sending error: ", err)
		return
	}
	log.Println("Sent error: ", message, "to", (*conn).RemoteAddr())
}
