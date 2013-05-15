package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

// http://www.pool.ntp.org/monitor/map

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

	// log.Println("RESP", serverMap["173.203.93.85"])
}

func getServerId(ip *net.IP) int {
	ipstr := ip.String()
	if status, ok := serverMap[ipstr]; ok {
		return status.Id
	}
	return 0
}
