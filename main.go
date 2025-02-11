package main

import (
	"log"

	"github.com/BhavyaMuni/protohackers/server"
)

func main() {
	// es := server.NewEchoServer()
	// go es.Start(":10000")

	log.Print("Starting server...")
	pts := server.NewPrimeTimeServer()
	go pts.Start(":10000")

	select {}
}
