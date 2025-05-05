package main

import (
	"log"

	"github.com/BhavyaMuni/protohackers/lineReversal"
	"github.com/BhavyaMuni/protohackers/server"
	"github.com/BhavyaMuni/protohackers/speedDaemon"
)

func main() {
	log.Print("Starting servers...")
	es := server.NewEchoServer()
	go es.Start(":10000")

	pts := server.NewPrimeTimeServer()
	go pts.Start(":10001")

	mtes := server.NewMeansToAnEndServer()
	go mtes.Start(":10002")

	bcs := server.NewBudgetChatServer()
	go bcs.Start(":10003")

	uds := server.NewUnusualDatabaseServer()
	go uds.Start(":10004")

	mims := server.NewMobInTheMiddleServer()
	go mims.Start(":10005")

	ssd := speedDaemon.NewSpeedDaemonServer()
	go ssd.Start(":10006")

	lrs := lineReversal.NewLineReversalServer()
	go lrs.Start(":10007")

	select {}
}
