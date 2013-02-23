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

func getServerData(id uint32) (logScores, error) {
	path := serverDataPath(id)
	url := "http://" + *s3host + "/servers/2012/" + path

	log.Println("getting URL", url)
	resp, err := http.Get(url)

	if err != nil || resp.StatusCode != 200 {
		log.Println("Could not get server data", resp.StatusCode, err)
		return nil, fmt.Errorf("%d %s", resp.StatusCode, err)
	}

	var reader *bufio.Reader
	reader = bufio.NewReader(resp.Body)
	defer resp.Body.Close()

	log.Println("Reading")

	scores := make(logScores, 0)

	for line, err := reader.ReadString('\n'); err == nil; line, err = reader.ReadString('\n') {
		// log.Println("getting line", line, err)
		score := new(logScore)
		err := json.Unmarshal([]byte(line), &score)
		if err != nil {
			return scores, err
		}
		scores = append(scores, score)

	}
	if err != nil && err != io.EOF {
		log.Println(err)
		return scores, err
	}

	log.Println("done")

	return scores, nil
}

func getRecentServerData(ip net.IP, since uint64, count int) (logScores, error) {
	url := fmt.Sprintf("http://%s/scores/%s/json?since=%d&limit=%d&monitor=*", *sitehost, ip.String(), since, count)
	log.Println("getting URL", url)

	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err != nil {
		log.Println("Could not get", url, err)
		return nil, fmt.Errorf("Could not get url %v: %v", url, err)
	}

	if err != nil || resp.StatusCode != 200 {
		log.Println("Could not get recent server data", resp.StatusCode, err)
		return nil, fmt.Errorf("Could not get recent server data: %d %s", resp.StatusCode, err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Could not read response", err)
		return nil, fmt.Errorf("Could not read response: %s", err)
	}

	data := historyData{}
	json.Unmarshal(body, &data)

	log.Println("Got number of recent scores:", len(data.History))

	return data.History, nil

}

func getMonitorData(ip net.IP, monitorChan chan serverMonitors) {
	resp, err := http.Get("http://" + *sitehost + "/scores/" + ip.String() + "/monitors?")
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != 200 {
		log.Println("Could not get monitor data", resp.StatusCode, err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Could not read response", err)
		return
	}

	monitors := make(map[string]serverMonitors)
	json.Unmarshal(body, &monitors)
	monitorData := monitors["monitors"]
	monitorChan <- monitorData
}
