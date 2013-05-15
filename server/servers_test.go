package server

import (
	. "launchpad.net/gocheck"
	"log"
	"testing"
)

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type ServerSuite struct {
}

var _ = Suite(&ServerSuite{})

func (s *ServerSuite) SetUpSuite(c *C) {
	log.SetFlags(log.Flags() + log.Lmicroseconds)
}

func (s *ServerSuite) TestPaths(c *C) {
	c.Check(NewServer(417).dataPath(), Equals, "400/417.json.gz")
	c.Check(NewServer(28).dataPath(), Equals, "0/28.json.gz")
}

func (s *ServerSuite) TestGetServerData(c *C) {
	srv := NewServer(101)
	scores, err := srv.GetData()
	c.Assert(err, IsNil)
	// log.Println("LEN", len(*scores))
	c.Assert(scores, HasLen, 44750)
	c.Assert(scores.First().Id, NotNil)
	c.Assert(scores.Last().Score, Equals, float64(19))
}

func (s *ServerSuite) TestGetRecentServerData(c *C) {
	limit := 4000
	srv := NewServer(101) // "207.171.7.151"
	scores, err := srv.getRecentData(1357768067, limit)
	c.Assert(err, IsNil)
	c.Assert(len(scores), Equals, limit)
	c.Assert(scores.Last().Id, NotNil)
	c.Assert(scores.First().Id, NotNil)
	//c.Assert(scores.Last().Score, Equals, float64(19))
}

func (s *ServerSuite) TestGetServerMonitors(c *C) {
	monitorChannel := make(chan serverMonitors)
	// "207.171.7.151"
	go getMonitorData(101, monitorChannel)
	monitors := <-monitorChannel
	c.Assert(monitors, HasLen, 1)
	c.Assert(monitors[0].Name, Equals, "Los Angeles, CA")
}
