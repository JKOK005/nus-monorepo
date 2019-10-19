package Evaluation

import (
	"fmt"
	"time"
	"context"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	dep "group-project/Utils"
	pb "group-project/Protobuf/Generate"
)

const (
	attempts = 10000  	// Put attempts
)

func generateRandVal() []byte {
	/*
	Generates a random value data pair.
	Value will be a UUID string cast into Bytes.
	*/
	if randUUID, err := uuid.NewRandom(); err != nil {
		panic(err)
	} else {
		return []byte(randUUID.String())
	}
}

func main() {
	// TODO: Do not let client connect with a defined URL:PORT. Create a service that queries ZK for the details
	bootstrap_url 	:= dep.GetEnvStr("REGISTER_LISTENER_DNS", "localhost")
	bootstrap_port 	:= uint32(dep.GetEnvInt("REGISTER_LISTENER_PORT", 8001))

	for key := 0; key < attempts; key++ {
		value := generateRandVal()
		glog.Infof(fmt.Sprintf("Attempting put key request - key: %s, val: %s"), string(key), string(value))

		if conn, err := grpc.Dial(fmt.Sprintf("%s:%d", bootstrap_url, bootstrap_port), grpc.WithInsecure()); err != nil {
			glog.Error(err)
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.PollTimeOutMs) * time.Millisecond)
			defer conn.Close(); defer cancel()

			client 		:= pb.NewPutKeyServiceClient(conn)
			resp, err 	:= client.PutKey(ctx, &pb.PutKeyMsg{Key: string(key), Val:value})
		}
	}
}