package server

import (
	"encoding/json"
	"fmt"
	"net"
)

type PrimeTimeRequest struct {
	Method string  `json:"method"`
	Number float64 `json:"number"`
}

type PrimeTimeResponse struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

type PrimeTimeServer struct {
	BaseServer
}

func (request PrimeTimeRequest) ValidRequest() bool {
	if request.Method != "isPrime" {
		return false
	}
	return true
}

func NewPrimeTimeServer() *PrimeTimeServer {
	s := &PrimeTimeServer{}
	s.HandleConnectionFunc = s.handleConnection
	return s
}

func (PrimeTimeServer) handleConnection(conn net.Conn) {
	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}

	var req PrimeTimeRequest
	err = json.Unmarshal(buf[:n], &req)
	if err != nil {
		fmt.Println("Error unmarshalling")
		return
	}

	fmt.Println(req)

	if !req.ValidRequest() {
		fmt.Println("Invalid")
		return
	}

	res := PrimeTimeResponse{Method: req.Method, Prime: IsPrime(req.Number)}

	data, _ := json.Marshal(res)
	conn.Write(data)
}

func IsPrime(n float64) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= int(n); i++ {
		if int(n)%i == 0 {
			return false
		}
	}
	return true
}
