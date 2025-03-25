package speedDaemon_test

import (
	"fmt"
	"testing"

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
	speed := speedDaemon.FindSpeed(1106, 10, 4869280, 4817702)
	if speed != 100 {
		t.Errorf("Expected speed to be 100, got %f", speed)
	}
	fmt.Println("Speed: ", speed)
}
