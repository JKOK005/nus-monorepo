package main

import (
	"flag"
	"fmt"
	"time"
	"os"
	"context"
	"strconv"
	"github.com/golang/glog"
	"go.uber.org/ratelimit"
	"google.golang.org/grpc"
	dep "group-project/Utils"
	pb "group-project/Protobuf/Generate"
	"sync"
	"math/rand"
)


func putKeys(attempts int, pollTimeOutMs int, bootstrap_url string,
			 bootstrap_ports []uint32, i int, name string) {

	file, err := os.Create(fmt.Sprint("Results/HashesAndThroughputIncreaseTest/",
									  name, "_", i, "?.txt"))
	if err != nil {
		panic(err)
	}

	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(),
								time.Duration(pollTimeOutMs) * time.Millisecond)
	defer cancel()

	var putClient = map[uint32]pb.PutKeyServiceClient{}
	var getClient = map[uint32]pb.GetKeyServiceClient{}
	for _, port := range(bootstrap_ports) {
		conn, _ := grpc.Dial(fmt.Sprintf("%s:%d", bootstrap_url, port),
							 grpc.WithInsecure())
		putClient[port] = pb.NewPutKeyServiceClient(conn)
		getClient[port] = pb.NewGetKeyServiceClient(conn)
		defer conn.Close()
	}

	file.WriteString("put,")

	// delay := time.Second
	// delayChange := delay / time.Duration(attempts)
	rl := ratelimit.New(10000)
	for key := 0; key < attempts; key++ {
		random_port := bootstrap_ports[rand.Intn(len(bootstrap_ports))]
		client := putClient[random_port]
		// start := time.Now()
		prev := time.Now()
		glog.Infof(fmt.Sprintf("Attempting put key request - key: %d, val: %d",
							   key, key))
		if resp, err := client.PutKey(ctx, &pb.PutKeyMsg{Key: strconv.Itoa(key),
								Val: []byte(strconv.Itoa(key))}); err != nil {
			panic(err)
		} else if resp.Ack != true {
			glog.Error("Failed to insert key: ", key)
		}
		now := rl.Take()
		// elapsed := time.Since(start)
		file.WriteString(fmt.Sprint(now.Sub(prev), ","))
		prev = now
		// time.Sleep(delay)
		// delay = delay - delayChange
	}

	file.WriteString("\nget,")

	// delay = time.Second
	for key := 0; key < attempts; key++ {
		random_port := bootstrap_ports[rand.Intn(len(bootstrap_ports))]
		client := getClient[random_port]
		// start := time.Now()
		prev := time.Now()
		glog.Infof(fmt.Sprintf("Attempting get key request - key: %d, val: %d",
							   key, key))
		if resp, err := client.GetKey(ctx,
						&pb.GetKeyMsg{Key: strconv.Itoa(key)}); err != nil {
			panic(err)
		} else if resp.Ack != true {
			glog.Error("Failed to retrieve key: ", key)
		} else {
			glog.Infof(fmt.Sprintf("Retrieved key: %d, val: %s", key,
									string(resp.Val)))
		}
		now := rl.Take()
		file.WriteString(fmt.Sprint(now.Sub(prev), ","))
		prev = now
		// elapsed := time.Since(start)
		// file.WriteString(fmt.Sprint(elapsed, ","))
		// time.Sleep(delay)
		// delay = delay - delayChange
	}
}

func main(){

	clients := flag.Int("clients", 1, "number of clients, should be int")
	servers := flag.Int("servers", 1, "number of servers, should be int")
	attempts := flag.Int("attempts", 100, "number of attempts, should be int")

	flag.Parse()

	ports := []int{}
	hashes := []int{}
	for i := 0; i < *servers; i ++ {
		ports = append(ports, 8000 + 2 * i)
		hashes = append(hashes, rand.Intn(10))
	}

	glog.Infof(fmt.Sprint("Used ports:", ports))
	glog.Infof(fmt.Sprint("Used hashes:", hashes))
	pollTimeOutMs := 600000000
	rand.Seed(time.Now().Unix())
	bootstrap_ports := []uint32{}
	bootstrap_url := dep.GetEnvStr("REGISTER_LISTENER_DNS", "localhost")


	var wg sync.WaitGroup

	for _, p := range(ports) {
		bootstrap_port := uint32(dep.GetEnvInt("REGISTER_LISTENER_PORT", p + 1))
		bootstrap_ports = append(bootstrap_ports, bootstrap_port)
	}

	for i := 0; i < *clients; i ++ {
		wg.Add(1)
		go putKeys(*attempts, pollTimeOutMs, bootstrap_url, bootstrap_ports, i,
				   fmt.Sprint(*clients, "_", *servers, "_", *attempts))
	}

	wg.Wait()
}
