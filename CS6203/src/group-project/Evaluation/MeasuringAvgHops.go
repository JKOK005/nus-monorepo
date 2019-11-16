package main

/*
	Experiment tests the average number of hops across the CHORD ring.

	As the number of hash group servers increase in the ring, we will expect an increase in random hops
	between 1 server to the rest.

	We conduct this experiment by bombarding a single server with random key PUT requests.
	We then measure the average hops taken by all requests.
	We expect for an essembly of N servers, the average hopes will follow a log(N) distribution due to the CHORD mechanism
*/

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	pb "group-project/Protobuf/Generate"
	dep "group-project/Utils"
	"strconv"
	"time"
)

func main() {
	// TODO: Do not let client connect with a defined URL:PORT. Create a service that queries ZK for the details
	// Leader address
	bootstrap_url 	:= dep.GetEnvStr("REGISTER_LISTENER_DNS", "localhost")
	bootstrap_port 	:= uint32(dep.GetEnvInt("REGISTER_LISTENER_PORT", 8001))

	const attempts 	= 1000
	pollTimeOutMs 	:= 100000000000000

	var totalHops int32 = 0
	/*
		Attempt to insert keys into a seed server
		- Client will be spawned to insert a range of keys into the leader
		- The leader is expected to route keys at random to the other severs
		- Client measures average number of hops
	*/
	if conn, err := grpc.Dial(fmt.Sprintf("%s:%d", bootstrap_url, bootstrap_port), grpc.WithInsecure()); err != nil {
		glog.Error(err)
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pollTimeOutMs) * time.Millisecond)
		client := pb.NewPutKeyServiceClient(conn)
		defer conn.Close(); defer cancel()

		for key := 0; key < attempts; key++ {
			glog.Infof(fmt.Sprintf("Attempting put key request - key: %d, val: %d", key, key))
			if resp, err := client.PutKey(ctx, &pb.PutKeyMsg{Key: strconv.Itoa(key),
				Val: []byte(strconv.Itoa(key))}); err != nil {
				panic(err)
			} else if resp.Ack != true {
				glog.Error("Failed to insert key: ", key)
			} else {
				totalHops += resp.Stats.NoOfHops
			}
		}
	}
	glog.Info("Hops: ", float64(totalHops) / float64(attempts))
}
