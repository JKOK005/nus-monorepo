package main

import (
	util "group-project/Utils"

	"flag"
	"fmt"
	"group-project/Services/Server"
	// "group-project/Services/Raft"
	"group-project/Services/Chord"
	dep "group-project/Utils"
	"sync"
	"time"
	"os"
	"strconv"
)

func testRocksDb() {
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

func main() {
	flag.Parse() // Needed for glog

	port64, _ := strconv.ParseUint(os.Args[3], 10, 32)
	baseHash64, _ := strconv.ParseUint(os.Args[2], 10, 32)
	port := int(port64)
	baseHash := uint32(baseHash64)

	nodeAddr := dep.GetEnvStr("REGISTER_LISTENER_DNS", "localhost")
	nodePort := uint32(dep.GetEnvInt("REGISTER_LISTENER_PORT", port))
	dbCli, _ := dep.InitRocksDB(dep.GetEnvStr("STORAGE_LOC", "./storage"))

	var wg sync.WaitGroup
	wg.Add(1)

	// Start up server to register all gRPC services
	go Server.Server{NodeAddr: nodeAddr, NodePort: nodePort,
					 DbCli: dbCli}.Start()

	// Start up state manager
	// go Raft.ElectionManager{NodeAddr: nodeAddr, NodePort: nodePort,
	// 						BaseHashGroup: baseHash, CycleNo: 0,
	// 						CyclesToTimeout: 10, CycleTimeMs: 1001,
	// 						State: Raft.Follower}.Start()

	time.Sleep(time.Second)

	go Chord.ChordManager{NodeAddr: nodeAddr, NodePort: nodePort,
						  BaseHashGroup: baseHash}.Start()


	time.Sleep(3 * time.Second)

	util.SetRequestChannel.ReqCh <-4
	time.Sleep(time.Second)

	util.SetPutChannel.ReqCh <-1
	time.Sleep(time.Second)

	wg.Wait()
}
