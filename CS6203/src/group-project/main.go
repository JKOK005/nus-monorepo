package main

import (
	"fmt"
	rocksdb "github.com/tecbot/gorocksdb"
)

func main ()  {
	bbto := rocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(rocksdb.NewLRUCache(3 << 30))
	opts := rocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	db, err := rocksdb.OpenDb(opts, "/Users/johan.kok/Desktop/rocksdb-test")

	if err != nil {
		panic(err)
	}

	ro := rocksdb.NewDefaultReadOptions()
	wo := rocksdb.NewDefaultWriteOptions()

	err = db.Put(wo, []byte("foo"), []byte("bar"))
	value, err := db.Get(ro, []byte("foo"))
	defer value.Free()

	fmt.Println("Value is: ", string(value.Data()))
}
