package store

import (
	"github.com/abh/dellingr/database"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
	"log"
	"time"
)

func (store *Store) getBucket() *s3.Bucket {
	auth := aws.Auth{store.User, store.Pass}

	region, ok := aws.Regions[store.BucketRegion]
	if !ok {
		log.Printf("Unknown region '%s'\n", store.BucketRegion)
		return nil
	}

	b := s3.New(auth, region).Bucket(store.Bucket)
	return b
}

func (store *Store) Update(serverId int) {
	db, err := database.New()
	if err != nil {
		log.Printf("Could not open database: '%s'\n", err)
		return
	}

	lastUpdate := db.GetLastUpdate(serverId)
	log.Println("lastUpdate", lastUpdate)
	if time.Since(lastUpdate) < time.Hour*24 {
		log.Printf("Server %d was updated recently, skipping", serverId)
		return
	}

	db.MarkUpdated(serverId, time.Now())

	log.Println("Updating server...")

	b := store.getBucket()
	b.GetReader(path)

}

// func (store *Store) List(serverId int) []string {

// 	b := store.getBucket()

// 	srv := server.NewServer(serverId)

// 	// path := srv.DataPath()
// 	// log.Println("PATH", path)

// 	rv := make([]string, 0)

// 	for _, year := range []string{"servers/2012/", "servers/2013/"} {

// 		list, err := b.List(year, "/", "", 100)
// 		if err != nil {
// 			log.Printf("s3 list error: '%s', %#v\n", err, err)
// 			return nil
// 		}

// 		log.Printf("Found %d keys\n", len(list.Contents))
// 		log.Printf("Common prefixes: %#v\n", list.CommonPrefixes)

// 		for _, f := range list.Contents {
// 			log.Printf(" file: %s %s\n", f.Key, f.StorageClass)
// 			rv = append(rv, f.Key)
// 		}

// 	}

// 	return rv

// 	// log.Printf("Got list: %#v\n", list)
// 	// log.Printf(" list contents: %#v\n", list.Contents)

// }
