package main

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	pb "group-project/Protobuf/Generate"
	dep "group-project/Utils"
	"time"
)

func main() {
	bootstrap_url 	:= dep.GetEnvStr("REGISTER_LISTENER_DNS", "localhost")
	bootstrap_port 	:= uint32(dep.GetEnvInt("REGISTER_LISTENER_PORT", 8001))

	pollTimeOutMs 	:= 10000

	glog.Infof(fmt.Sprintf("Attempting put request finger table recomputation"))

	if conn, err := grpc.Dial(fmt.Sprintf("%s:%d", bootstrap_url, bootstrap_port), grpc.WithInsecure()); err != nil {
		glog.Error(err)
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pollTimeOutMs) * time.Millisecond)
		defer conn.Close(); defer cancel()

		client := pb.NewRecomputeFingerTableServiceClient(conn)
		if resp, err := client.RecomputeFingerTable(ctx, &pb.RecomputeFingerTableMsg{}); err != nil {
			panic(err)
		} else if resp.Ack != true {
			glog.Error("Finger table recomputation failed")
		}
		glog.Info("Finger table recomputation succeeded")
	}
}