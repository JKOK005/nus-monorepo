package main

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

const (
	attempts = 500  	// Put attempts
)

func main() {
	// TODO: Do not let client connect with a defined URL:PORT. Create a service that queries ZK for the details
	// Leader address
	bootstrap_url 	:= dep.GetEnvStr("REGISTER_LISTENER_DNS", "localhost")
	bootstrap_port 	:= uint32(dep.GetEnvInt("REGISTER_LISTENER_PORT", 8001))

	// Slave address
	bootstrap_replica_url 	:= dep.GetEnvStr("REGISTER_LISTENER_SLAVE_DNS", "localhost")
	bootstrap_replica_port 	:= uint32(dep.GetEnvInt("REGISTER_LISTENER_SLAVE_PORT", 9001))

	pollTimeOutMs 	:= 10000

	/*
		Attempt to insert keys into leader and have keys replicate to slave
		- Client will be spawned to insert a range of keys into the leader
		- The leader will store all PUT requests and will replicate keys to slaves
		- Client will assert reads from slave

		Note: 	This setup requires 1 leader listening to in REGISTER_LISTENER_DNS:REGISTER_LISTENER_PORT
	 			and 1 slave listening to REGISTER_LISTENER_SLAVE_DNS:REGISTER_LISTENER_SLAVE_PORT
	*/
	for key := 0; key < attempts; key++ {
		glog.Infof(fmt.Sprintf("Attempting put key request - key: %d, val: %d", key, key))

		if conn, err := grpc.Dial(fmt.Sprintf("%s:%d", bootstrap_url, bootstrap_port), grpc.WithInsecure()); err != nil {
			glog.Error(err)
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pollTimeOutMs) * time.Millisecond)
			defer conn.Close(); defer cancel()

			client := pb.NewPutKeyServiceClient(conn)
			if resp, err := client.PutKey(ctx, &pb.PutKeyMsg{Key: strconv.Itoa(key),
															 Val: []byte(strconv.Itoa(key))}); err != nil {
				panic(err)
			} else if resp.Ack != true {
				glog.Error("Failed to insert key: ", key)
			}
		}
	}

	for key := 0; key < attempts; key++ {
		if conn, err := grpc.Dial(fmt.Sprintf("%s:%d", bootstrap_replica_url, bootstrap_replica_port), grpc.WithInsecure()); err != nil {
			glog.Error(err)
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pollTimeOutMs) * time.Millisecond)
			defer conn.Close(); defer cancel()

			client := pb.NewGetKeyServiceClient(conn)
			if resp, err := client.GetKey(ctx, &pb.GetKeyMsg{Key: strconv.Itoa(key)}); err != nil {
				panic(err)
			} else if resp.Ack != true {
				glog.Error("Failed to retrieve key: ", key)
			} else {
				glog.Infof(fmt.Sprintf("Retrieved key: %d, val: %s", key, string(resp.Val)))
			}
		}
	}
}