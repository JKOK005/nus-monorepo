package Raft

import (
	"fmt"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	pb "group-project/Protobuf/Generate"
	"group-project/Utils"
	"net"
)

type Server struct {
	nodeAddr 	string
	nodePort 	uint32
	termNo 		uint32
	dbCli 		*Utils.RocksDbClient
}

func (s *Server) setTermNo(newTermNo uint32) {s.termNo = newTermNo}

func (s *Server) Startup() error {
	/*
	Starts the RAFT GRPC server that listens to incoming requests
	*/
	glog.Infof("Starting server at %s:%d", s.nodeAddr, s.nodePort)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.nodeAddr, s.nodePort))
	if err != nil {return err}

	grpcServer := grpc.NewServer()
	pb.RegisterVotingServiceServer(grpcServer, &Server{})
	grpcServer.Serve(lis)

	return nil
}