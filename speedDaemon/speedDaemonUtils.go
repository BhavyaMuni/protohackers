package speedDaemon

import (
	"bufio"
	"encoding/binary"
	"errors"
	"log"
	"net"
)

const (
	ErrorMessageType         MessageType = 0x10
	PlateMessageType         MessageType = 0x20
	TicketMessageType        MessageType = 0x30
	WantHeartbeatMessageType MessageType = 0x40
	HeartbeatMessageType     MessageType = 0x41
	IAmCameraMessageType     MessageType = 0x80
	IAmDispatcherMessageType MessageType = 0x81
)

type MessageType byte

type Message interface {
	Handle(s *SpeedDaemonServer, conn *net.Conn)
}

type ClientType int

const (
	CAMERA ClientType = iota
	DISPATCHER
)

type ErrorMessage struct {
	MessageType
	Msg string
}

type TicketMessage struct {
	MessageType
	Plate      string
	Road       uint16
	Mile1      uint16
	Timestamp1 uint32
	Mile2      uint16
	Timestamp2 uint32
	Speed      uint16
}

func (m *TicketMessage) Handle(s *SpeedDaemonServer, conn *net.Conn) {
	log.Println("TicketMessage:", m)
}

func (m *IAmDispatcherMessage) Handle(s *SpeedDaemonServer, conn *net.Conn) {
	log.Println("IAmDispatcherMessage:", m)
	s.dispatchers[conn] = Dispatcher{NumRoads: m.NumRoads, Roads: m.Roads}
}

func ParseMessage(buf *bufio.Reader) (Message, MessageType, error) {
	log.Println("Parsing message")
	bufType, err := buf.Peek(1)
	if err != nil {
		return nil, ErrorMessageType, err
	}
	mType := MessageType(bufType[0])
	switch mType {
	case PlateMessageType:
		plateMsg := &PlateMessage{MessageType: mType}
		log.Println("PlateMessageType:", mType)
		buf.ReadByte()
		numLength, err := buf.ReadByte()
		if err != nil {
			return nil, PlateMessageType, err
		}
		plateBytes := make([]byte, numLength)
		_, err = buf.Read(plateBytes)
		if err != nil {
			return nil, PlateMessageType, err
		}
		plateMsg.Plate = string(plateBytes)
		err = binary.Read(buf, binary.BigEndian, &plateMsg.Timestamp)
		log.Println("PlateMessage:", plateMsg)
		if err != nil {
			return nil, PlateMessageType, err
		}
		return plateMsg, PlateMessageType, nil
	case TicketMessageType:
		ticketMsg := &TicketMessage{}
		err := binary.Read(buf, binary.BigEndian, ticketMsg)
		if err != nil {
			return nil, TicketMessageType, err
		}
		return ticketMsg, TicketMessageType, nil
	case WantHeartbeatMessageType:
		wantHeartbeatMsg := &WantHeartbeatMessage{}
		err := binary.Read(buf, binary.BigEndian, wantHeartbeatMsg)
		if err != nil {
			return nil, WantHeartbeatMessageType, err
		}
		return wantHeartbeatMsg, WantHeartbeatMessageType, nil
	case IAmCameraMessageType:
		iamCameraMsg := &IAmCameraMessage{}
		err := binary.Read(buf, binary.BigEndian, iamCameraMsg)
		if err != nil {
			return nil, IAmCameraMessageType, err
		}
		return iamCameraMsg, IAmCameraMessageType, nil
	case IAmDispatcherMessageType:
		iamDispatcherMsg := &IAmDispatcherMessage{}
		err := binary.Read(buf, binary.BigEndian, iamDispatcherMsg)
		if err != nil {
			return nil, IAmDispatcherMessageType, err
		}
		return iamDispatcherMsg, IAmDispatcherMessageType, nil
	default:
		return nil, ErrorMessageType, errors.New("unknown message type")
	}
}
