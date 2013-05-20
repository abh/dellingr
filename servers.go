package main

import (
	"encoding/json"
	"github.com/abh/dellingr/server"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

// http://www.pool.ntp.org/monitor/map

type Servers map[string]*server.Server

var servers Servers

var serverMap map[string]serverStatus

type serverStatus struct {
	Id      int  `json:"id"`
	Deleted bool `json:"deleted"`
}

func getServerMap() {
	log.Println("Getting server map")
	resp, err := http.Get("http://" + *sitehost + "/monitor/map")
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != 200 {
		log.Println("Could not get server map", resp.StatusCode, err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Could not read response", err)
		return
	}

	// log.Println("body", string(body))

	json.Unmarshal(body, &serverMap)

	log.Println("Got servermap")

	servers = make(Servers)

	for ipStr, status := range serverMap {
		ip := net.ParseIP(ipStr)
		srv := server.NewServer(status.Id)
		srv.Deleted = status.Deleted
		srv.Ip = &ip
		servers[ip.String()] = srv

	}

	// log.Println("RESP", serverMap["173.203.93.85"])
}

func getServer(ip *net.IP) *server.Server {
	// log.Printf("Servers: %#v\n", servers)
	if srv, ok := servers[ip.String()]; ok {
		return srv
	}
	return nil
}

func getServerId(ip *net.IP) int {
	if srv, ok := servers[ip.String()]; ok {
		return srv.Id
	}
	return 0
}
