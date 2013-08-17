package server

import (
	"fmt"
	"github.com/abh/dellingr/scores"
	. "launchpad.net/gocheck"
	"time"
)

type ScoreSuite struct {
	scores scores.LogScores
}

var _ = Suite(&ScoreSuite{})

func (s *ScoreSuite) SetUpSuite(c *C) {
	srv := testSrv()
	scores, err := srv.GetData(time.Unix(0, 0), time.Now())
	if err != nil {
		panic("Could not get server scores")
	}
	s.scores = scores.History
}

func (s *ScoreSuite) TestSample(c *C) {
	sampled := s.scores.Sample(22)
	points := 0
	for _, s := range sampled {
		c.Log(fmt.Sprintf("%#v", s))

		// the code will return more actual samples than requested
		// because it does one for each monitor (and the overall stats,
		// called "monitor 0")
		if s.MonitorId == 0 {
			points++
		}
	}

	c.Check(points, Equals, 22)
	// c.Log("Length: ", len(sampled))

	// sampled = s.scores.Sample(100)
	// c.Assert(sampled, HasLen, 100)
	// c.Log("Length: ", len(sampled))
}
