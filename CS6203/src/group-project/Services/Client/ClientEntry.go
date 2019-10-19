package Client

import (
	"fmt"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	pb "group-project/Protobuf/Generate"
	"net"
)

type Client struct {
	NodeAddr 	string
	NodePort 	uint32
}

func (c Client) Start() {
	/*
		Starts the RAFT GRPC server that listens to incoming requests
	*/
	glog.Infof("Starting client at %s:%d", c.NodeAddr, c.NodePort)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", c.NodeAddr, c.NodePort))
	if err != nil {glog.Fatal(err); panic(err)}

	grpcServer := grpc.NewServer()
	pb.RegisterGetKeyServiceServer(grpcServer, &c)
	pb.RegisterPutKeyServiceServer(grpcServer, &c)
	grpcServer.Serve(lis)
}