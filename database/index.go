package database

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

func (db *Database) getUpdated(idx string, serverId int) time.Time {
	c := db.Pool.Get()
	defer c.Close()

	lastUpdatedEpoch, err := redis.Int64(c.Do("ZSCORE", idx, serverId))
	if err != nil {
		log.Println("redis error", err)
		return time.Unix(0, 0)
	}
	if lastUpdatedEpoch == 0 {
		return time.Unix(0, 0)
	}
	log.Printf("lastUpdatedEpoch: %#v\n", lastUpdatedEpoch)

	return time.Unix(lastUpdatedEpoch, 0)
}

func (db *Database) setUpdated(idx string, serverId int, t time.Time) {
	c := db.Pool.Get()
	defer c.Close()

	ok, err := c.Do("ZADD", idx, t.Unix(), serverId)
	log.Println("Set says", ok, err)

	return
}

func (db *Database) GetFirstUpdated(serverId int) time.Time {
	return db.getUpdated("started", serverId)
}

func (db *Database) SetFirstUpdated(serverId int, t time.Time) {
	db.setUpdated("started", serverId, t)
}

func (db *Database) GetLastData(serverId int) time.Time {
	return db.getUpdated("last-data", serverId)
}

func (db *Database) SetLastData(serverId int, t time.Time) {
	db.setUpdated("last-data", serverId, t)
}

func (db *Database) GetLastUpdate(serverId int) time.Time {
	return db.getUpdated("updated", serverId)
}

func (db *Database) MarkUpdated(serverId int, t time.Time) {
	db.setUpdated("started", serverId, t)
}
