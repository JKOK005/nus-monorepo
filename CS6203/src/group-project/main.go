package main

import (
	"fmt"
	"group-project/Services/Raft"
	dep "group-project/Utils"
	"sync"
	"flag"
)

func testRocksDb ()  {
	rocksCli, err := dep.InitRocksDB("/Users/johan.kok/Desktop/rocksdb-test")

	if err != nil {
		panic(err)
	}

	// Test set and get key
	if err = rocksCli.Put("TestGetAndSetKey", []byte("TestGetAndSetVal")); err != nil {
		panic(err)
	} else {
		val, _ := rocksCli.Get("TestGetAndSetKey")
		fmt.Println("Value is: ", string(val))
	}

	// Test overwrite key
	if err = rocksCli.Put("TestGetAndSetKey", []byte("TestGetAndSetVal2")); err != nil {
		panic(err)
	} else {
		val, _ := rocksCli.Get("TestGetAndSetKey")
		fmt.Println("Value is: ", string(val))
	}

	// Test put immutable key
	if err = rocksCli.PutImmutable("TestGetAndSetImmutableKey", []byte("TestGetAndSetImmutableVal")); err != nil {
		panic(err)
	} else {
		val, _ := rocksCli.Get("TestGetAndSetImmutableKey")
		fmt.Println("Value is: ", string(val))
	}

	// Test rejection of key overwrite
	if err = rocksCli.PutImmutable("TestGetAndSetImmutableKey", []byte("TestGetAndSetImmutableVal")); err != nil {
		panic(err)
	}
}

func main () {
	flag.Parse()  // Needed for glog

	nodeAddr := "localhost"
	nodePort := uint32(8000)
	dbCli, _ := dep.InitRocksDB("./storage")

	var wg sync.WaitGroup
	wg.Add(1)

	// Start up server to register all gRPC services
	go Raft.Server{NodeAddr: nodeAddr, NodePort: nodePort, DbCli: dbCli}.Start()

	// Start up state manager
	go Raft.ElectionManager{ NodeAddr: nodeAddr, NodePort: nodePort, BaseHashGroup: 11, CycleNo: 0,
							 CyclesToTimeout: 10, CycleTimeMs: 1000, State: Raft.Follower}.Start()

	wg.Wait()
}