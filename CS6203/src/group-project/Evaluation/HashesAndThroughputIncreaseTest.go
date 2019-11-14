package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"group-project/Services/Client"
	"group-project/Services/DB"
	"group-project/Services/Election"
	"group-project/Services/Raft"
	"group-project/Services/Chord"
	dep "group-project/Utils"
	"math/rand"
	"sync"
)

func main(){
	flag.Parse()

	nodeAddr 		:= dep.GetEnvStr("REGISTER_LISTENER_DNS", "localhost")
	cycleNoStart 	:= uint32(dep.GetEnvInt("START_CYCLE_NO", 0))
	cyclesToTimeout := uint32(dep.GetEnvInt("CYCLES_TO_TIMEOUT", 10))
	startingState 	:= Election.Follower
	dbCli, _ 		:= dep.InitRocksDB(dep.GetEnvStr("STORAGE_LOC",
													 "./storage/leader"))

	ports := []int{8000, 8002, 8004, 8006, 8008}
	hashes := []int{2, 4, 6, 8, 10}
	glog.Infof(fmt.Sprint("Used ports:", ports))
	glog.Infof(fmt.Sprint("Used hashes:", hashes))

	var wg sync.WaitGroup

	for i, p := range(ports) {
		wg.Add(1)
<<<<<<< HEAD
		
=======

>>>>>>> 0d448efdf2ac819f37df5b5e47d57e8aa34a6296
		nodePort := uint32(dep.GetEnvInt("REGISTER_LISTENER_PORT", p))
		baseHashGroup := uint32(dep.GetEnvInt("HASH_GROUP", hashes[i]))
		cycleTimeMs := uint32(500 + rand.Intn(500))

		// Start up DB Client
		go DB.DbManager{DbCli:dbCli}.Start()

		// Register client services
		go Client.Client{NodeAddr: nodeAddr, NodePort: nodePort + 1}.Start()

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
		 					  BaseHashGroup: baseHashGroup,
							  FingerTable: nil, HighestHash: uint32(10)}.
							  Start()
	}

	wg.Wait()
}
