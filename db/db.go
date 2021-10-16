package db

import (
	"fmt"
	"os"

	"github.com/hhong0326/hhongcoin/utils"
	bolt "go.etcd.io/bbolt"
)

const (
	dbName       = "blockchain"
	dataBucket   = "data"
	blocksBucket = "blocks"
	checkpoint   = "checkpoint"
)

var db *bolt.DB

// 함수대신 struct가 adapter 역할
type DB struct{}

// storage interface method 구현
func (DB) FindBlock(hash string) []byte {
	return findBlock(hash)
}
func (DB) LoadChain() []byte {
	return loadChain()
}
func (DB) SaveBlock(hash string, data []byte) {
	saveBlock(hash, data)
}
func (DB) SaveChain(data []byte) {
	saveChain(data)
}
func (DB) EmptyBlocks() {
	emptyBlocks()
}

func getDBName() string {
	port := os.Args[2][7:]

	return fmt.Sprintf("%s_%s.db", dbName, port)
}

func InitDB() {
	if db == nil {
		dbPointer, err := bolt.Open(getDBName(), 0600, nil)
		db = dbPointer
		utils.HandleErr(err)
		// check or create
		err = db.Update(func(t *bolt.Tx) error {
			_, err := t.CreateBucketIfNotExists([]byte(dataBucket))
			_, err = t.CreateBucketIfNotExists([]byte(blocksBucket))
			return err
		})
		utils.HandleErr(err)

	}
}

// data 손상방지와 lock된 data 해제를 위해
func Close() {
	db.Close()
}

func saveBlock(hash string, data []byte) {
	err := db.Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data) // key, value
		return err
	})

	utils.HandleErr(err)
}

func saveChain(data []byte) {
	err := db.Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkpoint), data)
		return err
	})

	utils.HandleErr(err)
}

func loadChain() []byte {
	var data []byte
	db.View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})

	return data
}

func findBlock(hash string) []byte {
	var data []byte
	db.View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))

		return nil
	})

	return data
}

func emptyBlocks() {
	db.Update(func(t *bolt.Tx) error {
		utils.HandleErr(t.DeleteBucket([]byte(blocksBucket)))

		_, err := t.CreateBucket([]byte(blocksBucket))
		utils.HandleErr(err)
		return nil
	})
}
