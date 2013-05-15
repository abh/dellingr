package main

import (
	. "launchpad.net/gocheck"
	"log"
	"testing"
)

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type DellingrSuite struct {
}

var _ = Suite(&DellingrSuite{})

func (s *DellingrSuite) SetUpSuite(c *C) {
	log.SetFlags(log.Flags() + log.Lmicroseconds)
}

func (s *DellingrSuite) TestPaths(c *C) {
}
