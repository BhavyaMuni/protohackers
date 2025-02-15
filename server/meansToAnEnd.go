package server

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

type MeansToAnEndServer struct {
	BaseServer
}

type Message struct {
	Type byte
	P1   int32
	P2   int32
}

func NewMeansToAnEndServer() *MeansToAnEndServer {
	s := &MeansToAnEndServer{}
	s.HandleConnectionFunc = s.handleConnection
	return s
}

func findMean(queries map[int32]int32, minTime int32, maxTime int32) int32 {
	var sum int64
	var cnt int32
	for i := minTime; i <= maxTime; i++ {
		if v, ok := queries[i]; ok {
			sum += int64(v)
			cnt += 1
		}
	}

	if cnt == 0 {
		return 0
	}
	return int32(sum / int64(cnt))
}

func (MeansToAnEndServer) handleConnection(conn net.Conn) {
	log.Println("Connected with...")
	log.Println(int32(int64(581195615233) / int64(9438)))
	log.Println(conn.RemoteAddr())
	queries := make(map[int32]int32)
	for {
		message := Message{}
		err := binary.Read(conn, binary.BigEndian, &message)
		if err != nil {
			fmt.Println("binary.Read failed:", err)
			return
		}

		if message.Type == byte('I') {
			log.Println("Insert message received")
			queries[message.P1] = message.P2
		} else if message.Type == byte('Q') {
			log.Println("Query message received")
			mean := findMean(queries, message.P1, message.P2)
			binary.Write(conn, binary.BigEndian, mean)
		} else {
			binary.Write(conn, binary.BigEndian, "Invalid message received")
			log.Println("Invalid message received")
		}

		log.Printf("Received message: %v", message)
	}
}
