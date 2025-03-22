package speedDaemon

import "net"

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

func (m *IAmDispatcherMessage) Handle(s *SpeedDaemonServer, conn *net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, road := range m.Roads {
		s.roadDispatchers[road] = Dispatcher{NumRoads: m.NumRoads, Roads: m.Roads, Conn: conn}
	}
}
