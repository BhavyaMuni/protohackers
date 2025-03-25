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
	if m.Interval <= 0 {
		log.Println("Heartbeat interval is 0, not sending heartbeat")
		return
	}
	go SendHeartbeat(conn, m.Interval)
}

func SendHeartbeat(conn *net.Conn, interval uint32) {
	tickerDuration := time.Duration(float64(interval)) * time.Millisecond * 100
	log.Println("Heartbeat interval: ", tickerDuration)
	if tickerDuration <= 0 {
		log.Println("Heartbeat interval is 0, not sending heartbeat")
		return
	}
	ticker := time.NewTicker(tickerDuration)
	for range ticker.C {
		heartbeatMsg := HeartbeatMessage{MessageType: HeartbeatMessageType}
		err := binary.Write(*conn, binary.BigEndian, heartbeatMsg)
		if err != nil {
			log.Println("Error sending heartbeat message: ", err)
			return
		}
	}
}
