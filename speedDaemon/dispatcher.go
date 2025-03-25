package speedDaemon

import (
	"encoding/binary"
	"log"
	"net"
)

type IAmDispatcherMessage struct {
	MessageType
	NumRoads uint8
	Roads    []uint16
}

type Dispatcher struct {
	NumRoads uint8
	Roads    []uint16
	Conn     *net.Conn
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

func (m *IAmDispatcherMessage) Handle(s *SpeedDaemonServer, conn *net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, road := range m.Roads {
		s.roadDispatchers[road] = Dispatcher{NumRoads: m.NumRoads, Roads: m.Roads, Conn: conn}
	}
}

func (d *Dispatcher) CheckSpeedViolation(observations []Observation, currentObservation Observation) {
	for i := range len(observations) - 1 {
		if observations[i].Camera.Road == currentObservation.Camera.Road {
			distance := currentObservation.Camera.Mile - observations[i].Camera.Mile
			time := currentObservation.Timestamp - observations[i].Timestamp
			speed := uint16(uint32(distance) / time)
			if speed > currentObservation.Camera.Limit {
				if currentObservation.Timestamp > observations[i].Timestamp {
					d.SendTicket(observations[i], currentObservation, speed)
				} else {
					d.SendTicket(currentObservation, observations[i], speed)
				}
			}
		}
	}
}

func (d *Dispatcher) SendTicket(observation1 Observation, observation2 Observation, speed uint16) {
	ticketMsg := TicketMessage{
		MessageType: TicketMessageType,
		Road:        observation1.Camera.Road,
		Mile1:       observation1.Camera.Mile,
		Timestamp1:  observation1.Timestamp,
		Mile2:       observation2.Camera.Mile,
		Timestamp2:  observation2.Timestamp,
		Speed:       speed,
	}
	plateBytes := []byte(observation1.Plate)
	numPlateBytes := uint8(len(plateBytes))

	plateBytes = append(plateBytes, numPlateBytes)
	ticketMsg.Plate = string(plateBytes)

	err := binary.Write(*d.Conn, binary.BigEndian, ticketMsg)
	if err != nil {
		log.Println("Error sending ticket:", err)
	}
}
