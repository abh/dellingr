package scores

import (
	"fmt"
	. "launchpad.net/gocheck"
	"testing"
	"time"
)

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type ScoreSuite struct {
	scores LogScores
}

var _ = Suite(&ScoreSuite{})

func (s *ScoreSuite) SetUpSuite(c *C) {
	s.scores = LogScores{}

	now := time.Unix(time.Now().Unix()-(2000*60*15), 0)

	// Simulate data
	for i := int(0); i < 5000; i++ {
		ls := new(LogScore)
		ls.Id = int64(i)
		ls.MonitorId = i % 3
		ls.Offset = i / 1000
		ls.Score = 19.5
		ls.Step = 5
		ls.Ts = uint64(now.Unix())
		s.scores = append(s.scores, ls)
		now.Add(time.Minute * 15)
	}
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
