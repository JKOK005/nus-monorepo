package Server

import (
	"fmt"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	pb "group-project/Protobuf/Generate"
	"group-project/Utils"
	"net"
)

type Server struct {
	NodeAddr 	string
	NodePort 	uint32
	DbCli 		*Utils.RocksDbClient
}

func (s Server) Start() {
	/*
	Starts the RAFT GRPC server that listens to incoming requests
	*/
	glog.Infof("Starting server at %s:%d", s.NodeAddr, s.NodePort)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.NodeAddr, s.NodePort))
	if err != nil {glog.Fatal(err); panic(err)}

	grpcServer := grpc.NewServer()
	pb.RegisterVotingServiceServer(grpcServer, &s)
	pb.RegisterHeartbeatServiceServer(grpcServer, &s)
	grpcServer.Serve(lis)
}
