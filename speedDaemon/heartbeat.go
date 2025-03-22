package speedDaemon

import (
	"encoding/binary"
	"log"
	"net"
	"time"
)

type HeartbeatMessage struct {
	MessageType
}

type WantHeartbeatMessage struct {
	MessageType
	Interval uint32
}

func (m *WantHeartbeatMessage) Handle(s *SpeedDaemonServer, conn *net.Conn) {
	go SendHeartbeat(conn, m.Interval)
}

func SendHeartbeat(conn *net.Conn, interval uint32) {
	ticker := time.NewTicker(time.Duration(interval/10) * time.Second)
	for range ticker.C {
		heartbeatMsg := HeartbeatMessage{MessageType: HeartbeatMessageType}
		err := binary.Write(*conn, binary.BigEndian, heartbeatMsg)
		if err != nil {
			log.Println("Error sending heartbeat message: ", err)
		}
	}
}
