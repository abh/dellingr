package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(2)

	log.New(os.Stderr, "dellingr", log.LstdFlags)
	log.SetPrefix("dellingr ")
	log.SetFlags(log.Flags() + log.Lmicroseconds)

	log.Println("Starting")

	getServerMap()

	go httpServer()

	terminate := make(chan os.Signal)
	signal.Notify(terminate, os.Interrupt)

	<-terminate
	log.Printf("signal received, stopping")
}
