package speedDaemon

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"sync"
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
			go dispatcher.MonitorTicketQueue(s.tickets[road], &s.ticketDays, &s.mu)
		} else {
			go dispatcher.MonitorTicketQueue(s.tickets[road], &s.ticketDays, &s.mu)
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

func (d *Dispatcher) MonitorTicketQueue(tickets <-chan *Ticket, ticketDays *map[uint32]map[string]bool, mu *sync.Mutex) {
	for ticket := range tickets {
		mu.Lock()
		day := ticket.Timestamp1 / 86400
		if _, ok := (*ticketDays)[day]; !ok {
			(*ticketDays)[day] = make(map[string]bool)
		}
		if _, ok := (*ticketDays)[day][ticket.Plate]; !ok {
			go d.SendTicket(*ticket)
		}
		(*ticketDays)[day][ticket.Plate] = true
		log.Println("Ticket days: ", (*ticketDays)[day])
		mu.Unlock()
	}
}
