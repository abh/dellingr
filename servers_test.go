package main

import (
	. "launchpad.net/gocheck"
	"log"
	"net"
	"testing"
)

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type TestSuite struct {
}

var _ = Suite(&TestSuite{})

func (s *TestSuite) TestPaths(c *C) {

	log.SetFlags(log.Flags() + log.Lmicroseconds)

	c.Check(serverDataPath(417), Equals, "400/417.json.gz")
	c.Check(serverDataPath(28), Equals, "0/28.json.gz")
}

func (s *TestSuite) TestGetServerData(c *C) {
	scores, err := getServerData(101)
	c.Assert(err, IsNil)
	// log.Println("LEN", len(*scores))
	c.Assert(*scores, HasLen, 44750)
	c.Assert(scores.First().Id, NotNil)
	c.Assert(scores.Last().Score, Equals, float64(19))
}

func (s *TestSuite) TestGetRecentServerData(c *C) {
	limit := 4000
	scores, err := getRecentServerData(net.ParseIP("207.171.7.151"), 1357768067, limit)
	c.Assert(err, IsNil)
	c.Assert(len(*scores), Equals, limit)
	c.Assert(scores.Last().Id, NotNil)
	c.Assert(scores.First().Id, NotNil)
	//c.Assert(scores.Last().Score, Equals, float64(19))
}

func (s *TestSuite) TestGetServerMonitors(c *C) {

	monitorChannel := make(chan serverMonitors)
	go getMonitorData(net.ParseIP("207.171.7.151"), monitorChannel)

	monitors := <-monitorChannel

	c.Assert(monitors[0].Name, Equals, "Los Angeles, CA")
}
