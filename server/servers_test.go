package server

import (
	. "launchpad.net/gocheck"
	"log"
	"net"
	"testing"
)

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type ServerSuite struct {
}

var _ = Suite(&ServerSuite{})

func testSrv() *Server {
	srv := NewServer(101)
	srv.Ip = net.ParseIP("207.171.7.151")
	return srv
}

func (s *ServerSuite) SetUpSuite(c *C) {
	log.SetFlags(log.Flags() + log.Lmicroseconds)
}

func (s *ServerSuite) TestPaths(c *C) {
	c.Check(NewServer(417).dataPath(), Equals, "400/417.json.gz")
	c.Check(NewServer(28).dataPath(), Equals, "0/28.json.gz")
}

func (s *ServerSuite) TestGetServerData(c *C) {
	srv := testSrv()
	scores, err := srv.GetAllData()
	c.Assert(err, IsNil)
	// c.Assert(scores.History, HasLen, 44750)
	c.Assert(scores.History.First().Id, NotNil)
	// c.Assert(scores.History.Last().Score, Equals, float64(19))
}

func (s *ServerSuite) TestGetRecentServerData(c *C) {
	limit := 4000
	srv := testSrv()
	scores, err := srv.getRecentData(1357768067, limit)
	c.Assert(err, IsNil)
	c.Assert(len(scores), Equals, limit)
	c.Assert(scores.Last().Id, NotNil)
	c.Assert(scores.First().Id, NotNil)
	//c.Assert(scores.Last().Score, Equals, float64(19))
}

func (s *ServerSuite) TestGetServerMonitors(c *C) {
	monitorChannel := make(chan serverMonitors)
	srv := testSrv()
	go srv.getMonitorData(monitorChannel)
	monitors := <-monitorChannel
	hasMonitors := len(monitors) > 0
	c.Assert(hasMonitors, Equals, true)
	c.Assert(monitors[0].Name, Equals, "Los Angeles, CA")
}
