package main

import (
	"net"

	"github.com/BhavyaMuni/protohackers/server"
)

func main() {
	es := server.NewEchoServer()
	go es.Start(":10000")

	pts := server.NewPrimeTimeServer()
	go pts.Start(":10001")

	select {}
}

type Server interface {
	HandleConnection(conn net.Conn) // Handles a single connection
	Start(port string) error        // Starts the server on the specified port
}

// func ConvertBytesToStruct(res []byte) PrimeTimeRequest {
// 	var request PrimeTimeRequest
// 	err := json.Unmarshal(res, &request)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return request
// }

// func ConvertStructToBytes(res PrimeTimeResponse) []byte {
// 	data, _ := json.Marshal(res)
// 	return data
// }
