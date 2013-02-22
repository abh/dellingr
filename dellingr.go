package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
)

var (
	sitehost = flag.String("site", "www.beta.grundclock.com", "hostname for NTP Pool site")
	s3host   = flag.String("s3", "s3beta.ntppool.org", "Hostname for web server with data files")
)

func main() {
	runtime.GOMAXPROCS(2)

	flag.Parse()

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
