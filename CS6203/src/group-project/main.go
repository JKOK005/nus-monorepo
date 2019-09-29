package main

import (
	"fmt"
	"group-project/Services/Raft"
	dep "group-project/Utils"
	"math/rand"
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

	nodeAddr 		:= dep.GetEnvStr("REGISTER_LISTENER_DNS", "localhost")
	nodePort 		:= uint32(dep.GetEnvInt("REGISTER_LISTENER_PORT", 8000))
	baseHashGroup 	:= uint32(dep.GetEnvInt("HASH_GROUP", 1))
	cycleNoStart 	:= uint32(dep.GetEnvInt("START_CYCLE_NO", 0))
	cyclesToTimeout := uint32(dep.GetEnvInt("CYCLES_TO_TIMEOUT", 10))
	cycleTimeMs 	:= uint32(500 + rand.Intn(500)) // Generates a random value between 0.5 - 1 sec
	startingState 	:= Raft.Follower
	dbCli, _ 		:= dep.InitRocksDB(dep.GetEnvStr("STORAGE_LOC", "./storage"))

	var wg sync.WaitGroup
	wg.Add(1)

	// Start up server to register all gRPC services
	go Raft.Server{NodeAddr: nodeAddr, NodePort: nodePort, DbCli: dbCli}.Start()

	// Start up state manager
	go Raft.ElectionManager{ NodeAddr: nodeAddr, NodePort: nodePort, BaseHashGroup: baseHashGroup, CycleNo: cycleNoStart,
							 CyclesToTimeout: cyclesToTimeout, CycleTimeMs: cycleTimeMs, State: startingState}.Start()

	wg.Wait()
}