package main

import (
	"flag"
	"github.com/abh/dellingr/server"
	"github.com/abh/dellingr/store"
	"log"
	"os"
	"os/signal"
	"runtime"
)

var (
	sitehost = flag.String("site", "www.beta.grundclock.com", "hostname for NTP Pool site")
	s3host   = flag.String("s3", "s3beta.ntppool.org", "Hostname for web server with data files")
	s3update = flag.Bool("s3update", false, "Update S3 and exit")
)

func init() {
	runtime.GOMAXPROCS(2)

	flag.Parse()

	server.SiteHost = *sitehost
	server.S3Host = *s3host

}

func main() {
	log.New(os.Stderr, "dellingr", log.LstdFlags)
	log.SetPrefix("dellingr ")
	log.SetFlags(log.Flags() + log.Lmicroseconds)

	getServerMap()

	if *s3update {
		log.Println("Updating S3")

		s3key := os.Getenv("S3KEY")
		s3secret := os.Getenv("S3SECRET")
		s3bucket := os.Getenv("S3BUCKET")
		s3region := os.Getenv("S3REGION")

		store := store.New(s3key, s3secret, s3bucket, s3region)
		for ip, server := range serverMap {
			log.Printf("Updating %s\n", ip)
			store.Update(server.Id)
		}

		return
	}

	log.Println("Starting http server")

	go httpServer()

	terminate := make(chan os.Signal)
	signal.Notify(terminate, os.Interrupt)

	<-terminate
	log.Printf("signal received, stopping")
}
