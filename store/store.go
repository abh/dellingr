package store

type Store struct {
	User         string
	Pass         string
	Bucket       string
	BucketRegion string
}

func New(user, pass, bucket, bucketRegion string) *Store {
	store := new(Store)
	store.User = user
	store.Pass = pass
	store.Bucket = bucket
	store.BucketRegion = bucketRegion
	return store
}
