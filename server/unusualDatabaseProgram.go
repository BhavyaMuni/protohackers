package server

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type UnusualDatabaseServer struct {
	BaseServer
	Database map[string]string
}

func NewUnusualDatabaseServer() *UnusualDatabaseServer {
	uds := &UnusualDatabaseServer{}
	uds.Database = make(map[string]string)
	uds.Database["version"] = "1.0"
	return uds
}

func (s *UnusualDatabaseServer) Start(port string) error {
	udpAddr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return err
	}

	// Create a UDP listener
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Error listening on UDP:", err)
		return err
	}
	defer conn.Close()

	fmt.Println("UDP server listening on", udpAddr)

	buffer := make([]byte, 1024)

	log.Println("Database:", s.Database)

	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading UDP packet:", err)
			continue
		}
		log.Println("Received request from", addr, ":", string(buffer[:n]))
		received := string(buffer[:n])
		response := s.handleDatabaseRequest(strings.TrimSpace(received))
		conn.WriteTo([]byte(response), addr)
		log.Println("Sent response to", addr, ":", response)
	}
}

func (s *UnusualDatabaseServer) handleDatabaseRequest(request string) string {
	if strings.Contains(request, "=") {
		parts := strings.SplitN(request, "=", 2)
		key := parts[0]
		value := ""
		if len(parts) > 1 {
			value = parts[1]
		}
		if key == "version" {
			return ""
		}
		s.Database[key] = value
		return ""
	} else {
		if value, ok := s.Database[request]; ok {
			return request + "=" + value
		}
		return ""
	}
}
