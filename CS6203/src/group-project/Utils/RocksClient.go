package Utils

import (
	"errors"
	"fmt"
	rocksdb "github.com/tecbot/gorocksdb"
)

type RocksDbClient struct {
	dbConn *rocksdb.DB
}

func InitRocksDB(path string) (*RocksDbClient, error) {
	/*
		Initializes a rocks DB connection

		@Input:
		path - rocks DB path to file system, for storage
	*/
	bbto := rocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(rocksdb.NewLRUCache(3 << 30))
	opts := rocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)

	if db, err := rocksdb.OpenDb(opts, path); err != nil {
		return nil, err
	} else {
		return &RocksDbClient{dbConn: db}, nil
	}
}

func (r *RocksDbClient) Get(key string) ([]byte, error) {
	/*
		Retrieves key from rocks DB connection

		@Input:
		key - key to retrieve
	*/
	ro := rocksdb.NewDefaultReadOptions()
	if val, err := r.dbConn.Get(ro, []byte(key)); err != nil {
		return nil, err
	} else {
		return val.Data(), nil
	}
}

func (r *RocksDbClient) Put(key string, val []byte) error {
	/*
		Sets key from rocks DB connection

		@Input:
		key - key to insert
		val - value to insert
	*/
	wo := rocksdb.NewDefaultWriteOptions()
	if err := r.dbConn.Put(wo, []byte(key), val); err != nil {
		return err
	} else {
		return nil
	}
}

func (r *RocksDbClient) PutImmutable(key string, val []byte) error {
	/*
		Sets key from rocks DB connection. First asserts immutability of keys.

		@Input:
		key - key to insert
		val - value to insert
	*/
	if firstVal, err := r.Get(key); err != nil {
		return err
	} else {
		if len(firstVal) == 0 {
			return r.Put(key, val)
		} else {
			return errors.New(fmt.Sprintf("Key %s exists. Insertion forbidden.", key))
		}
	}
}