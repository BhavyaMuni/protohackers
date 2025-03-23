package speedDaemon_test

import (
	"encoding/binary"
	"log"
	"net"
	"testing"
	"time"

	"github.com/BhavyaMuni/protohackers/speedDaemon"
)

type TestCameraMessage struct {
	MessageType byte
	Road        uint16
}

type TestPlateMessage struct {
	MessageType byte
	NumPlates   uint8
	Plates      [6]byte
	Timestamp   uint32
}

func TestCamera(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:10006")
	if err != nil {
		log.Fatal(err)
	}
	iamCameraMsg := &speedDaemon.IAmCameraMessage{
		MessageType: speedDaemon.IAmCameraMessageType,
		Road:        1,
	}
	binary.Write(conn, binary.BigEndian, iamCameraMsg)
	plateMsg := &TestPlateMessage{
		MessageType: 32,
		NumPlates:   6,
		Plates:      [6]byte{'A', 'B', 'C', '1', '2', '3'},
		Timestamp:   uint32(time.Now().Unix()),
	}
	err = binary.Write(conn, binary.BigEndian, plateMsg)
	if err != nil {
		log.Fatal(err)
	}
}
