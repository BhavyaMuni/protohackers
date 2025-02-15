package main

import (
	"log"

	"github.com/BhavyaMuni/protohackers/server"
)

func main() {
	log.Print("Starting servers...")
	es := server.NewEchoServer()
	go es.Start(":10000")

	pts := server.NewPrimeTimeServer()
	go pts.Start(":10001")

	mtes := server.NewMeansToAnEndServer()
	go mtes.Start(":10002")
	select {}
}
