package speedDaemon

import (
	"bytes"
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

type Ticket struct {
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
		dispatcher := Dispatcher{NumRoads: m.NumRoads, Roads: m.Roads, Conn: conn}
		if _, ok := s.tickets[road]; !ok {
			s.tickets[road] = make(chan *Ticket)
			go dispatcher.MonitorTicketQueue(s, s.tickets[road])
		} else {
			go dispatcher.MonitorTicketQueue(s, s.tickets[road])
		}
	}
}

func (d *Dispatcher) SendTicket(ticket Ticket) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, TicketMessageType)
	binary.Write(buf, binary.BigEndian, uint8(len(ticket.Plate)))
	buf.WriteString(ticket.Plate)
	binary.Write(buf, binary.BigEndian, ticket.Road)
	binary.Write(buf, binary.BigEndian, ticket.Mile1)
	binary.Write(buf, binary.BigEndian, ticket.Timestamp1)
	binary.Write(buf, binary.BigEndian, ticket.Mile2)
	binary.Write(buf, binary.BigEndian, ticket.Timestamp2)
	binary.Write(buf, binary.BigEndian, ticket.Speed)
	binary.Write(*d.Conn, binary.BigEndian, buf.Bytes())
	log.Println("Sent ticket: ", ticket)
}

func (d *Dispatcher) MonitorTicketQueue(server *SpeedDaemonServer, tickets <-chan *Ticket) {
	for ticket := range tickets {
		server.mu.Lock()
		day1 := ticket.Timestamp1 / 86400
		day2 := ticket.Timestamp2 / 86400
		if _, ok := server.ticketDays[day1]; !ok {
			server.ticketDays[day1] = make(map[string]bool)
		}
		if _, ok := server.ticketDays[day2]; !ok {
			server.ticketDays[day2] = make(map[string]bool)
		}
		_, d1ok := server.ticketDays[day1][ticket.Plate]
		_, d2ok := server.ticketDays[day2][ticket.Plate]
		if !d1ok && !d2ok {
			go d.SendTicket(*ticket)
		} else {
			log.Println("Ticket already sent: ", ticket.Plate)
		}
		server.ticketDays[day1][ticket.Plate] = true
		server.ticketDays[day2][ticket.Plate] = true
		server.mu.Unlock()
	}
}
