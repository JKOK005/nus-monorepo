package Raft

import (
	"fmt"
	"github.com/golang/glog"
	pb "group-project/Protobuf/Generate"
	"group-project/Utils"
	"net"
	"google.golang.org/grpc"
)

type Server struct {
	nodeAddr string;
	nodePort uint32;
	dbClient Utils.RocksDbClient;
	termNo uint32;
}

func (s *Server) setTermNo(newTermNo uint32) {s.termNo = newTermNo}

func (s *Server) preLoadCommitIndex() error {
	/*
	Sets commit index to 0 if commit index key does not exist. This assumes that the DB is newly created.
	If key exists, we do not overwrite the previous value of commit index.
	*/
	if keyExists, err := s.dbClient.Exists("commit_msg_indx"); err != nil {
		return err
	} else {
		if !keyExists {
			_ = s.dbClient.Put("commit_msg_indx", []byte("0"))
		}
		return nil
	}
}

func (s *Server) Startup() error {
	/*
	Starts the RAFT GRPC server that listens to incoming requests
	Also initializes commit message pointer for logs to 0 if message pointer does not exist
	*/
	if err := s.preLoadCommitIndex(); err != nil {return err}

	glog.Infof("Starting server at %s:%d", s.nodeAddr, s.nodePort)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.nodeAddr, s.nodePort))
	if err != nil {return err}

	grpcServer := grpc.NewServer()
	pb.RegisterVotingServiceServer(grpcServer, &Server{})
	grpcServer.Serve(lis)

	return nil
}