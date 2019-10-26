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
	// "os"
	// "strconv"
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

	// port64, _ := strconv.ParseUint(os.Args[1], 10, 32)
	// baseHash64, _ := strconv.ParseUint(os.Args[2], 10, 32)
	// port := int(port64)
	// baseHash := uint32(baseHash64)

	nodeAddr := dep.GetEnvStr("REGISTER_LISTENER_DNS", "localhost")
	nodePort1 := uint32(dep.GetEnvInt("REGISTER_LISTENER_PORT", 8001))
	nodePort2 := uint32(dep.GetEnvInt("REGISTER_LISTENER_PORT", 8002))
	nodePort3 := uint32(dep.GetEnvInt("REGISTER_LISTENER_PORT", 8003))
	nodePort4 := uint32(dep.GetEnvInt("REGISTER_LISTENER_PORT", 8004))
	dbCli, _ := dep.InitRocksDB(dep.GetEnvStr("STORAGE_LOC", "./storage"))

	var wg sync.WaitGroup
	wg.Add(1)

	// Start up server to register all gRPC services
	go Server.Server{NodeAddr: nodeAddr, NodePort: nodePort1,
					 DbCli: dbCli}.Start()
	go Server.Server{NodeAddr: nodeAddr, NodePort: nodePort2,
					 DbCli: dbCli}.Start()
	go Server.Server{NodeAddr: nodeAddr, NodePort: nodePort3,
					 DbCli: dbCli}.Start()
	go Server.Server{NodeAddr: nodeAddr, NodePort: nodePort4,
					 DbCli: dbCli}.Start()

	// Start up state manager
	// go Raft.ElectionManager{NodeAddr: nodeAddr, NodePort: nodePort, BaseHashGroup: 11, CycleNo: 0,
	// 	CyclesToTimeout: 10, CycleTimeMs: 1001, State: Raft.Follower}.Start()

	time.Sleep(time.Second)

	go Chord.QueryManager{NodeAddr: nodeAddr, NodePort: nodePort1,
						  BaseHashGroup: 1}.Start()
	go Chord.QueryManager{NodeAddr: nodeAddr, NodePort: nodePort2,
  						  BaseHashGroup: 4}.Start()
	go Chord.QueryManager{NodeAddr: nodeAddr, NodePort: nodePort3,
						  BaseHashGroup: 7}.Start()
	go Chord.QueryManager{NodeAddr: nodeAddr, NodePort: nodePort4,
				  		  BaseHashGroup: 10}.Start()


	time.Sleep(3 * time.Second)

	// util.SetRequestChannel.ReqCh <-10
	// time.Sleep(time.Second)

	// util.SetRequestChannel.ReqCh <-7
	// time.Sleep(time.Second)
	//
	// util.SetRequestChannel.ReqCh <-4
	// time.Sleep(time.Second)
	//
	// util.SetRequestChannel.ReqCh <-1
	// time.Sleep(time.Second)
	//
	// util.SetRequestChannel.ReqCh <-10
	// time.Sleep(time.Second)
	//
	// util.SetRequestChannel.ReqCh <-7
	// time.Sleep(time.Second)
	//
	// util.SetRequestChannel.ReqCh <-4
	// time.Sleep(time.Second)
	//
	// util.SetRequestChannel.ReqCh <-1
	// time.Sleep(time.Second)

	util.SetPutChannel.ReqCh <-1
	time.Sleep(time.Second)

	wg.Wait()
}
