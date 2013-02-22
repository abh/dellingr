package main

import (
	"bufio"
	// "compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

// http://www.pool.ntp.org/monitor/map

var serverMap map[string]serverStatus

type serverMonitors []serverMonitor

type serverMonitor struct {
	Id    uint32  `json:"id"`
	Name  string  `json:"name"`
	Score float64 `json:"score"`
}

type serverStatus struct {
	Id      uint32 `json:"id"`
	Deleted bool   `json:"deleted"`
}

func getServerMap() {
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

	// log.Println("RESP", serverMap["173.203.93.85"])
}

func getServerId(ip *net.IP) uint32 {
	ipstr := ip.String()
	if status, ok := serverMap[ipstr]; ok {
		return status.Id
	}
	return 0
}

func serverDataPath(id uint32) string {
	hundreds := id - id%100
	return fmt.Sprintf("%v/%v.json.gz", hundreds, id)
}

func getServerData(id uint32) (*logScores, error) {
	path := serverDataPath(id)
	// url := "http://localhost/~ask/ntp/" + path

	// url := "http://10.0.201.231:6081/servers/2012/" + path
	// url := "http://ntpbeta.s3.amazonaws.com/servers/2012/" + path

	url := "http://" + *s3host + "/servers/2012/" + path

	resp, err := http.Get(url)

	log.Println("getting URL", url)

	scores := make(logScores, 0)
	// scores := make([]interface{}, 0)

	if err != nil {
		log.Println("Could not get", url, err)
		return &scores, fmt.Errorf("Could not get url %v: %v", url, err)
	}

	var reader *bufio.Reader

	// gzReader, err := gzip.NewReader(resp.Body)
	// if err == nil {
	// log.Println("reading gzip")
	// reader = bufio.NewReader(gzReader)
	// } else {
	// gzReader.Close()
	log.Println("not gzip")
	reader = bufio.NewReader(resp.Body)
	// }
	defer resp.Body.Close()

	log.Println("Reading")

	for line, err := reader.ReadString('\n'); err == nil; line, err = reader.ReadString('\n') {
		// log.Println("getting line", line, err)
		score := new(logScore)
		err := json.Unmarshal([]byte(line), &score)
		if err != nil {
			return &scores, err
		}
		scores = append(scores, score)

	}
	if err != nil && err != io.EOF {
		log.Println(err)
		return &scores, err
	}

	log.Println("done")

	return &scores, nil

}
