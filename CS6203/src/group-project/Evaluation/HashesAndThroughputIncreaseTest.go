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


func putKeys(pollTimeOutMs int, putClient map[uint32]pb.PutKeyServiceClient,
			 bootstrap_ports []uint32, i int, file *os.File,
			 ctx context.Context) {

	time1 := time.Now()

	key := rand.Intn(10)
	random_port := bootstrap_ports[rand.Intn(len(bootstrap_ports))]
	client := putClient[random_port]

	glog.Infof(fmt.Sprintf("Attempting put key request - key: %d, val: %d",
						   key, key))
	if resp, err := client.PutKey(ctx, &pb.PutKeyMsg{Key: strconv.Itoa(key),
							Val: []byte(strconv.Itoa(key))}); err != nil {
		panic(err)
	} else if resp.Ack != true {
		glog.Error("Failed to insert key: ", key)
	}
	time2 := time.Now()
	file.WriteString(fmt.Sprint(time2.Sub(time1).Nanoseconds() / 1000000, ","))
}

func main(){

	clients := flag.Int("clients", 1, "number of clients, should be int")
	servers := flag.Int("servers", 1, "number of servers, should be int")
	rate := flag.Int("rate", 1, "ratelimit, should be int")

	flag.Parse()

	ports := []int{}
	for i := 0; i < *servers; i ++ {
		ports = append(ports, 8000 + 2 * i)
	}

	glog.Infof(fmt.Sprint("Used ports:", ports))
	pollTimeOutMs := 600000000
	rand.Seed(time.Now().Unix())
	bootstrap_ports := []uint32{}
	bootstrap_url := dep.GetEnvStr("REGISTER_LISTENER_DNS", "localhost")
	name := fmt.Sprint(*clients, "_", *rate, "_", *servers)

	var wg sync.WaitGroup

	for _, p := range(ports) {
		bootstrap_port := uint32(dep.GetEnvInt("REGISTER_LISTENER_PORT", p + 1))
		bootstrap_ports = append(bootstrap_ports, bootstrap_port)
	}

	wg.Add(1)

	file, err := os.Create(fmt.Sprint("Results/HashesAndThroughputIncreaseTest/",
									  name , ".txt"))
	if err != nil {
		panic(err)
	}
	defer file.Close() // type = os.File


	ctx, cancel := context.WithTimeout(context.Background(),
								time.Duration(pollTimeOutMs) * time.Millisecond)
	defer cancel()

	var putClient = map[uint32]pb.PutKeyServiceClient{}
	for _, port := range(bootstrap_ports) {
		conn, _ := grpc.Dial(fmt.Sprintf("%s:%d", bootstrap_url, port),
							 grpc.WithInsecure())
		putClient[port] = pb.NewPutKeyServiceClient(conn)
		defer conn.Close()
	}


	rl := ratelimit.New(*rate) // type = ratelimit.Limiter
	for i := 0; i < *clients; i ++ {
		rl.Take()
		go putKeys(pollTimeOutMs, putClient, bootstrap_ports, i, file, ctx)
	}

	wg.Wait()
}
