package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/abh/dellingr/scores"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	SiteHost = "www.beta.grundclock.com"
	S3Host   = "s3beta.ntppool.org"
)

type serverMonitors []serverMonitor

type Server struct {
	Id int
}

type serverMonitor struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Score float64 `json:"score"`
}

type HistoryData struct {
	History  scores.LogScores  `json:"history"`
	Monitors serverMonitors    `json:"monitors"`
	Server   map[string]string `json:"server"`
}

func NewServer(id int) *Server {
	return &Server{Id: id}
}

func (s *Server) DataPath() string {
	hundreds := s.Id - s.Id%100
	return fmt.Sprintf("%v/%v.json.gz", hundreds, s.Id)
}

func (s *Server) GetData() (*HistoryData, error) {
	monitorChannel := make(chan serverMonitors)
	go getMonitorData(s.Id, monitorChannel)

	scores, err := s.getArchiveData()
	if err != nil {
		log.Println("Error fetching server data", err)
	}

	// 	scores := &logScores{}

	since := uint64(0)
	if scores != nil && len(scores) > 0 {
		if lastScore := scores.Last(); lastScore.Ts > 0 {
			// TODO(abh) Round the Ts back to midnight so pagination cache better
			since = lastScore.Ts
		}
	}

	for {
		recentScores, err := s.getRecentData(since, 4000)
		if err != nil {
			return nil, err
		}
		if len(recentScores) > 0 {
			since = recentScores.Last().Ts
			// log.Println("Got recent scores", len(*recentScores), len(*scores))
			if scores == nil {
				scores = recentScores
			} else {
				scores = append(scores, recentScores...)
			}

			log.Println("new length", len(scores))
		} else {
			log.Println("didn't get recent scores!")
			break
		}
	}

	monitors := <-monitorChannel

	history := &HistoryData{}
	history.Monitors = monitors

	return history, nil
}

func (s *Server) getArchiveData() (scores.LogScores, error) {
	path := s.DataPath()
	url := "http://" + S3Host + "/servers/2012/" + path

	log.Println("getting URL", url)
	resp, err := http.Get(url)

	if err != nil || resp.StatusCode != 200 {
		status := 500
		if resp != nil {
			status = resp.StatusCode
		}
		log.Println("Could not get server data", status, err)
		return nil, fmt.Errorf("%d %s", status, err)
	}

	var reader *bufio.Reader
	reader = bufio.NewReader(resp.Body)
	defer resp.Body.Close()

	log.Println("Reading")

	srvScores := make(scores.LogScores, 0)

	for line, err := reader.ReadString('\n'); err == nil; line, err = reader.ReadString('\n') {
		// log.Println("getting line", line, err)
		score := new(scores.LogScore)
		err := json.Unmarshal([]byte(line), &score)
		if err != nil {
			return srvScores, err
		}
		srvScores = append(srvScores, score)

	}
	if err != nil && err != io.EOF {
		log.Println(err)
		return srvScores, err
	}

	log.Println("done")

	return srvScores, nil
}

func (s *Server) getRecentData(since uint64, count int) (scores.LogScores, error) {
	url := fmt.Sprintf("http://%s/scores/%s/json?since=%d&limit=%d&monitor=*", SiteHost, s.Id, since, count)
	log.Println("getting URL", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Could not get", url, err)
		return nil, fmt.Errorf("Could not get url %v: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("Could not get recent server data", resp.StatusCode, err)
		return nil, fmt.Errorf("Could not get recent server data: %d %s", resp.StatusCode, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Could not read response", err)
		return nil, fmt.Errorf("Could not read response: %s", err)
	}

	data := HistoryData{}
	json.Unmarshal(body, &data)

	log.Println("Got number of recent scores:", len(data.History))

	return data.History, nil

}

func getMonitorData(id int, monitorChan chan serverMonitors) {
	resp, err := http.Get(fmt.Sprintf("http://%s/scores/%d/monitors", SiteHost, id))
	if err != nil || resp.StatusCode != 200 {
		log.Println("Could not get monitor data", err)
		if resp != nil {
			log.Println("Monitor data response code", resp.StatusCode)
		}
		monitorChan <- make(serverMonitors, 0)
		return
	}
	defer resp.Body.Close()
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
