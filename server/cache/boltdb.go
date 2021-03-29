package cache

import (
	"github.com/segmentio/encoding/json"
	bolt "go.etcd.io/bbolt"
	"os"
)

var DefaultBucket = []byte("default_bucket")

type BoltDb struct {
	*bolt.DB
	bucket []byte
}

func NewBoltDb(path string) *BoltDb {
	db, err := bolt.Open(path, os.ModePerm, nil)
	if err != nil {
		panic(err) //TODO: Probably shouldn't do this
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(DefaultBucket)
		return err
	})
	if err != nil {
		panic(err) //TODO: Probably shouldn't do this
	}
	return &BoltDb{DB: db}
}

func (db *BoltDb) Get(key string, val interface{}) {
	_ = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(DefaultBucket)
		valBytes := b.Get([]byte(key))
		if valBytes != nil {
			_ = json.Unmarshal(valBytes, val)
		}

		return nil
	})
}

func (db *BoltDb) Put(key string, val interface{}) {
	_ = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(DefaultBucket)
		valBytes, _ := json.Marshal(val)
		_ = b.Put([]byte(key), valBytes)

		return nil
	})
}
