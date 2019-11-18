package main

import (
	"flag"
	"fmt"
	"group-project/Services/Client"
	"group-project/Services/DB"
	"group-project/Services/Election"
	"group-project/Services/Raft"
	"group-project/Services/Chord"
	dep "group-project/Utils"
	"math/rand"
	"sync"
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

	port := flag.Int("port", 8000, "the port of the server, should be an int")
	hash := flag.Int("hash", 1, "the hash of the server, should be an int")
	storage := flag.String("storage", "1", "the location of the storage")

	flag.Parse()  // Needed for glog

	nodeAddr 		:= dep.GetEnvStr("REGISTER_LISTENER_DNS", "localhost")
	nodePort 		:= uint32(dep.GetEnvInt("REGISTER_LISTENER_PORT", *port))
	baseHashGroup 	:= uint32(dep.GetEnvInt("HASH_GROUP", *hash))
	cycleNoStart 	:= uint32(dep.GetEnvInt("START_CYCLE_NO", 0))
	cyclesToTimeout := uint32(dep.GetEnvInt("CYCLES_TO_TIMEOUT", 10))
	cycleTimeMs 	:= uint32(500 + rand.Intn(500)) // Generates a random value between 0.5 - 1 sec
	startingState 	:= Election.Follower
	dbCli, _ 		:= dep.InitRocksDB(dep.GetEnvStr("STORAGE_LOC",
									   fmt.Sprintf("./storage/%s", *storage)))

	// fmt.Print(fmt.Sprint("\n\n\n\n\n\n\n", fmt.Sprintf("./storage/%s", *storage), "\n\n\n\n\n\n\n"))

	var wg sync.WaitGroup
	wg.Add(1)

	// Start up DB Client
	go DB.DbManager{DbCli:dbCli}.Start()

	// Register client services
	go Client.Client{NodeAddr: nodeAddr, NodePort: nodePort +1}.Start()

	// Start up server to register all gRPC services
	go Raft.Server{NodeAddr: nodeAddr, NodePort: nodePort}.Start()

	// Start up state manager
	go Election.ElectionManager{NodeAddr: nodeAddr, NodePort: nodePort,
								BaseHashGroup: baseHashGroup,
								CycleNo: cycleNoStart,
							 	CyclesToTimeout: cyclesToTimeout,
								CycleTimeMs: cycleTimeMs,
								State: startingState}.Start()

	// Start chord manager
	go Chord.ChordManager{NodeAddr: nodeAddr, NodePort: nodePort,
	 					  BaseHashGroup: baseHashGroup, FingerTable: nil,
						  HighestHash: uint32(10)}.Start()

	wg.Wait()
}
