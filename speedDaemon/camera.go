package speedDaemon

import (
	"log"
	"net"
)

type IAmCameraMessage struct {
	MessageType
	Road  uint16
	Mile  uint16
	Limit uint16
}

type Camera struct {
	Road  uint16
	Mile  uint16
	Limit uint16
	Conn  *net.Conn
}

type PlateMessage struct {
	MessageType
	Plate     string
	Timestamp uint32
}

type Observation struct {
	Plate     string
	Timestamp uint32
	Camera    Camera
}

func (m *PlateMessage) Handle(s *SpeedDaemonServer, conn *net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.observations[m.Plate]; !ok {
		s.observations[m.Plate] = []Observation{}
	}

	newObservation := Observation{
		Plate:     m.Plate,
		Timestamp: m.Timestamp,
		Camera:    s.cameras[conn],
	}
	dispatcher := s.roadDispatchers[newObservation.Camera.Road]
	go dispatcher.CheckSpeedViolation(s.observations[m.Plate], newObservation)
	s.observations[m.Plate] = append(s.observations[m.Plate], newObservation)
}

func (m *IAmCameraMessage) Handle(s *SpeedDaemonServer, conn *net.Conn) {
	log.Println("IAmCameraMessage:", m)
	s.cameras[conn] = Camera{Road: m.Road, Mile: m.Mile, Limit: m.Limit, Conn: conn}
}
