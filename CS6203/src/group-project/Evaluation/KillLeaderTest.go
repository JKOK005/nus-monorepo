package main

import (
	"flag"
	"fmt"
	"time"
	// "os"
	// "context"
	// "strconv"
	"github.com/golang/glog"
	// "google.golang.org/grpc"
	"group-project/Services/Client"
	"group-project/Services/DB"
	"group-project/Services/Election"
	"group-project/Services/Raft"
	"group-project/Services/Chord"
	dep "group-project/Utils"
	// pb "group-project/Protobuf/Generate"
	"math/rand"
	"sync"
)


func putKeys(attempts int, pollTimeOutMs int, bootstrap_url string,
			 bootstrap_port uint32, bootstrap_replica_url string,
			 bootstrap_replica_port uint32) {

	glog.Info(fmt.Sprint(attempts, pollTimeOutMs, bootstrap_url, bootstrap_port,
						 bootstrap_replica_url, bootstrap_replica_port))

	// file, err := os.Create(fmt.Sprint("Results/HashesAndThroughputIncreaseTest/",
	// 								  bootstrap_port))
	// if err != nil {
	// 	panic(err)
	// }
	//
	// defer file.Close()

	// start := time.Now()

	// if conn, err := grpc.Dial(fmt.Sprintf("%s:%d", bootstrap_url,
	// 						bootstrap_port), grpc.WithInsecure()); err != nil {
	// 		glog.Error(err)
	// } else {
	// 	ctx, cancel := context.WithTimeout(context.Background(),
	// 							time.Duration(pollTimeOutMs) * time.Millisecond)
	// 	client := pb.NewPutKeyServiceClient(conn)
	// 	defer conn.Close(); defer cancel()
	//
	// 	for key := 0; key < attempts; key++ {
	// 		glog.Infof(fmt.Sprintf("Attempting put key request - key: %d, val: %d", key, key))
	// 		if resp, err := client.PutKey(ctx, &pb.PutKeyMsg{Key: strconv.Itoa(key),
	// 			Val: []byte(strconv.Itoa(key))}); err != nil {
	// 			panic(err)
	// 		} else if resp.Ack != true {
	// 			glog.Error("Failed to insert key: ", key)
	// 		}
	// 	}
	// }
	//
	// // time.Sleep(1 * time.Second)
	//
	// if conn, err := grpc.Dial(fmt.Sprintf("%s:%d", bootstrap_replica_url,
	// 				bootstrap_replica_port), grpc.WithInsecure()); err != nil {
	// 	glog.Error(err)
	// } else {
	// 	ctx, cancel := context.WithTimeout(context.Background(),
	// 							time.Duration(pollTimeOutMs) * time.Millisecond)
	// 	client := pb.NewGetKeyServiceClient(conn)
	// 	defer conn.Close(); defer cancel()
	//
	// 	for key := 0; key < attempts; key++ {
	// 		if resp, err := client.GetKey(ctx,
	// 						&pb.GetKeyMsg{Key: strconv.Itoa(key)}); err != nil {
	// 			panic(err)
	// 		} else if resp.Ack != true {
	// 			glog.Error("Failed to retrieve key: ", key)
	// 		} else {
	// 			glog.Infof(fmt.Sprintf("Retrieved key: %d, val: %s", key,
	// 									string(resp.Val)))
	// 		}
	// 	}
	// }

	// elapsed := time.Since(start)
	// n3, _ := file.WriteString(fmt.Sprint(elapsed, ","))
    // fmt.Printf("wrote %d bytes\n", n3)
}

func main(){
	flag.Parse()

	nodeAddr := dep.GetEnvStr("REGISTER_LISTENER_DNS", "localhost")
	cycleNoStart := uint32(dep.GetEnvInt("START_CYCLE_NO", 0))
	cyclesToTimeout := uint32(dep.GetEnvInt("CYCLES_TO_TIMEOUT", 10))
	startingState := Election.Follower
	dbCli, _ := dep.InitRocksDB(dep.GetEnvStr("STORAGE_LOC","./storage/leader"))

	ports := []int{8000, 8002, 8004, 8006, 8008}
	hashes := []int{1, 1, 1, 1, 1}
	glog.Infof(fmt.Sprint("Used ports:", ports))
	glog.Infof(fmt.Sprint("Used hashes:", hashes))
	pollTimeOutMs := 6000000

	var wg sync.WaitGroup
	quitChannel := make(map[int](chan bool))

	for i, p := range(ports) {

		quitChannel[p] = make(chan bool)

		wg.Add(1)

		attempts := 100
		nodePort := uint32(dep.GetEnvInt("REGISTER_LISTENER_PORT", p))
		baseHashGroup := uint32(dep.GetEnvInt("HASH_GROUP", hashes[i]))
		cycleTimeMs := uint32(500 + rand.Intn(500))

		// Start up DB Client
		go DB.DbManager{DbCli:dbCli}.Start()

		// Register client services
		go Client.Client{NodeAddr: nodeAddr, NodePort: nodePort + 1}.Start()

		// Start up server to register all gRPC services
		go func(){
			switch {
			case condition:

			}

			Raft.Server{NodeAddr: nodeAddr, NodePort: nodePort}.Start()
		}


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
							  FingerTable: nil, HighestHash: uint32(10)}.Start()

		time.Sleep(4 * time.Second)

		bootstrap_url := dep.GetEnvStr("REGISTER_LISTENER_DNS", "localhost")
		bootstrap_port := uint32(dep.GetEnvInt("REGISTER_LISTENER_PORT", p + 1))
		bootstrap_replica_url := dep.GetEnvStr("REGISTER_LISTENER_SLAVE_DNS",
											   "localhost")
		bootstrap_replica_port := uint32(dep.GetEnvInt("REGISTER_LISTENER_SLAVE_PORT",
													   p + 1001))

		go putKeys(attempts, pollTimeOutMs, bootstrap_url, bootstrap_port,
				   bootstrap_replica_url, bootstrap_replica_port)
	}

	wg.Wait()
}
