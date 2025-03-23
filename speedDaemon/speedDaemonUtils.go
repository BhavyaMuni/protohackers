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
		if err != nil {
			return nil, PlateMessageType, err
		}
		return plateMsg, PlateMessageType, nil
	case WantHeartbeatMessageType:
		wantHeartbeatMsg := &WantHeartbeatMessage{}
		err := binary.Read(buf, binary.BigEndian, wantHeartbeatMsg)
		if err != nil {
			return nil, WantHeartbeatMessageType, err
		}
		return wantHeartbeatMsg, WantHeartbeatMessageType, nil
	case IAmCameraMessageType:
		iAmCameraMsg := &IAmCameraMessage{}
		err := binary.Read(buf, binary.BigEndian, iAmCameraMsg)
		if err != nil {
			return nil, IAmCameraMessageType, err
		}
		return iAmCameraMsg, IAmCameraMessageType, nil
	case IAmDispatcherMessageType:
		iAmDispatcherMsg := &IAmDispatcherMessage{
			MessageType: IAmDispatcherMessageType,
			Roads:       make([]uint16, 0),
		}
		buf.ReadByte()
		numRoads, err := buf.ReadByte()
		if err != nil {
			return nil, IAmDispatcherMessageType, err
		}
		iAmDispatcherMsg.NumRoads = numRoads
		roads := make([]uint16, numRoads)
		err = binary.Read(buf, binary.BigEndian, roads)
		if err != nil {
			return nil, IAmDispatcherMessageType, err
		}
		iAmDispatcherMsg.Roads = roads
		return iAmDispatcherMsg, IAmDispatcherMessageType, nil
	default:
		return nil, ErrorMessageType, errors.New("unknown message type")
	}
}
