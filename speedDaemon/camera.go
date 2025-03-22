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
	// go calculateSpeed(s.observations[m.Plate], newObservation)
	s.observations[m.Plate] = append(s.observations[m.Plate], newObservation)
}

func (m *IAmCameraMessage) Handle(s *SpeedDaemonServer, conn *net.Conn) {
	log.Println("IAmCameraMessage:", m)
	s.cameras[conn] = Camera{Road: m.Road, Mile: m.Mile, Limit: m.Limit, Conn: conn}
}

func calculateSpeed(observations []Observation, currentObservation Observation) {
	for i := 0; i < len(observations)-1; i++ {
		if observations[i].Camera.Road == currentObservation.Camera.Road {
			distance := currentObservation.Camera.Mile - observations[i].Camera.Mile
			time := currentObservation.Timestamp - observations[i].Timestamp
			speed := uint32(distance) / time
			log.Printf("Speed: %d, Distance: %d, Time: %d\n", speed, distance, time)
		}
	}
}
