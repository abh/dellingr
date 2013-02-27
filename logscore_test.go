package main

import (
	. "launchpad.net/gocheck"
)

type ScoreSuite struct {
	scores logScores
}

var _ = Suite(&ScoreSuite{})

func (s *ScoreSuite) SetUpSuite(c *C) {
	scores, err := getServerData(101)
	if err != nil {
		panic("Could not get server scores")
	}
	s.scores = scores
}

func (s *ScoreSuite) TestSample(c *C) {
	sampled := s.scores.Sample(20)
	points := 0
	for _, s := range sampled {
		if s.MonitorId == 0 {
			points++
		}
	}

	c.Assert(points, Equals, 20)
	// c.Log("Length: ", len(sampled))

	// sampled = s.scores.Sample(100)
	// c.Assert(sampled, HasLen, 100)
	// c.Log("Length: ", len(sampled))
}
