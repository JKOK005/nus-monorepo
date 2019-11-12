package main

/*
	Experiment demonstrates the trade off between redundancy and client PUT request timing.

	In the present model, we adopt a full replication protocol in which the leader must establish
	the PUT across all slaves in order to guarantee the client that the write has been successful.
	This clearly implies that the more slaves there are listening to the leader, the longer the expected
	wait time to replicate data across all slaves
*/

import (
	"fmt"
	"context"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	dep "group-project/Utils"
	pb "group-project/Protobuf/Generate"
	"strconv"
	"time"
)

func main() {
	// TODO: Do not let client connect with a defined URL:PORT. Create a service that queries ZK for the details
	// Leader address
	bootstrap_url := dep.GetEnvStr("REGISTER_LISTENER_DNS", "localhost")
	bootstrap_port := uint32(dep.GetEnvInt("REGISTER_LISTENER_PORT", 8001))

	const attempts int 	= 10000
	pollTimeOutMs 		:= 10000
	var totalTime int64 = 0

	// Repeat experiment with 1 / 2 ... N slaves listening to the leader
	if conn, err := grpc.Dial(fmt.Sprintf("%s:%d", bootstrap_url, bootstrap_port), grpc.WithInsecure()); err != nil {
		glog.Error(err)
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pollTimeOutMs) * time.Millisecond)
		defer conn.Close(); defer cancel()

		for key := 0; key < attempts; key++ {
			glog.Infof(fmt.Sprintf("Attempting put key request - key: %d, val: %d", key, key))
			startTime := time.Now()
			client := pb.NewPutKeyServiceClient(conn)
			if resp, err := client.PutKey(ctx, &pb.PutKeyMsg{Key: strconv.Itoa(key),
				Val: []byte(strconv.Itoa(key))}); err != nil {
				panic(err)
			} else if resp.Ack != true {
				glog.Error("Failed to insert key: ", key)
			}
			elapsedTime := time.Since(startTime)
			totalTime += elapsedTime.Nanoseconds()
		}

	}
	glog.Info("Average insertion time: ", totalTime / int64(attempts))
}