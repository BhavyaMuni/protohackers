package speedDaemon_test

import (
	"bufio"
	"encoding/binary"
	"log"
	"net"
	"testing"

	"github.com/BhavyaMuni/protohackers/speedDaemon"
)

type TestDispatcherMessage struct {
	MessageType byte
	NumRoads    uint8
	Roads       [1]uint16
}

func TestDispatcher(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:10006")
	if err != nil {
		log.Fatal(err)
	}
	iamDispatcherMsg := &TestDispatcherMessage{
		MessageType: 129,
		NumRoads:    1,
		Roads:       [1]uint16{66},
	}
	//convert to binary data
	err = binary.Write(conn, binary.BigEndian, iamDispatcherMsg)
	if err != nil {
		log.Fatal(err)
	}

	buf := bufio.NewReader(conn)
	for {
		message, messageType, err := speedDaemon.ParseMessage(buf)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(message)
		log.Println(messageType)
	}
}
