package main

import (
	. "launchpad.net/gocheck"
	"os"
)

type S3Suite struct {
}

var _ = Suite(&S3Suite{})

func (s *S3Suite) SetUpSuite(c *C) {
}

func (s *S3Suite) TestS3(c *C) {

	s3key := os.Getenv("S3KEY")
	s3secret := os.Getenv("S3SECRET")
	s3bucket := os.Getenv("S3BUCKET")
	s3region := os.Getenv("S3REGION")

	store := NewStore(s3key, s3secret, s3bucket, s3region)
	c.Assert(store, NotNil)
	c.Log("store.BucketRegion", store.BucketRegion)
	store.Get(177)

}
