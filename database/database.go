// Copyright 2013 Ask Bjørn Hansen
// Copyright 2012 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Redis keys and types:
//
// maxPackageId string: next id to assign
// id:<path> string: id for given import path
// pkg:<id> hash
//      terms: space separated search terms
//      path: import path
//      synopsis: synopsis
//      gob: snappy compressed gob encoded doc.Package
//      rank: document search rank
//      etag:
//      kind: p=package, c=command, d=directory with no go files
// index:<term> set: package ids for given search term
// index:import:<path> set: packages with import path
// index:project:<root> set: packages in project with root
// crawl zset: package id, Unix time for next crawl
// block set: packages to block

// Package database manages storage for GoPkgDoc.
package database

import (
	"flag"
	"log"
	"net/url"
	"os"
	"time"

	// "code.google.com/p/snappy-go/snappy"
	"github.com/garyburd/redigo/redis"
)

type Database struct {
	Pool interface {
		Get() redis.Conn
	}
}

type Package struct {
	Path     string `json:"path"`
	Synopsis string `json:"synopsis,omitempty"`
}

type byPath []Package

func (p byPath) Len() int           { return len(p) }
func (p byPath) Less(i, j int) bool { return p[i].Path < p[j].Path }
func (p byPath) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

var (
	redisServer      = flag.String("db-server", "redis://127.0.0.1:6379", "URI of Redis server.")
	redisIdleTimeout = flag.Duration("db-idle-timeout", 250*time.Second, "Close Redis connections after remaining idle for this duration.")
	redisLog         = flag.Bool("db-log", false, "Log database commands")
)

func dialDb() (c redis.Conn, err error) {
	u, err := url.Parse(*redisServer)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil && c != nil {
			c.Close()
		}
	}()

	c, err = redis.Dial("tcp", u.Host)
	if err != nil {
		return
	}

	if *redisLog {
		l := log.New(os.Stderr, "", log.LstdFlags)
		c = redis.NewLoggingConn(c, l, "")
	}

	if u.User != nil {
		if pw, ok := u.User.Password(); ok {
			if _, err = c.Do("AUTH", pw); err != nil {
				return
			}
		}
	}
	return
}

// New creates a database configured from command line flags.
func New() (*Database, error) {
	pool := &redis.Pool{
		Dial:        dialDb,
		MaxIdle:     10,
		IdleTimeout: *redisIdleTimeout,
	}

	if c := pool.Get(); c.Err() != nil {
		return nil, c.Err()
	} else {
		c.Close()
	}

	return &Database{Pool: pool}, nil
}

// Exists returns true if package with import path exists in the database.
func (db *Database) Exists(path string) (bool, error) {
	c := db.Pool.Get()
	defer c.Close()
	return redis.Bool(c.Do("EXISTS", "id:"+path))
}
