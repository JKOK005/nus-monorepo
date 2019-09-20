package Raft

import (
	"context"
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
}

var (
	termNo uint32 = 0;
)

func setTermNo(newTermNo uint32) {termNo = newTermNo}

func (s *Server) RequestVote(ctx context.Context, msg *pb.RequestForVoteMsg) (*pb.RequestForVoteReply, error) {
	/*
	Evaluates if node should vote positive for a vote request RPC.

	Vote is YES if the incoming vote message V contains:
		- V.candidateTerm > V.node_candidateTerm
		- V.logLength >= V.node_logLength

	Vote is NO if otherwise
	*/
	indexAsBytes, _ := s.dbClient.Get("commit_msg_indx")
	nodeLogLength 	:= uint32(indexAsBytes[0])
	if termNo < msg.CandidateTerm && nodeLogLength <= msg.LogLength {
		glog.Infof("%s:%d - node term no: %d < %d & logLength: %d <= %d",
							s.nodeAddr, s.nodePort, termNo, msg.CandidateTerm, nodeLogLength, msg.LogLength)
		glog.Info("Vote acknowledged")
		setTermNo(msg.CandidateTerm)
		return &pb.RequestForVoteReply{Ack:true}, nil
	} else {
		glog.Infof("%s:%d - node term no: %d >= %d & logLength: %d > %d",
			s.nodeAddr, s.nodePort, termNo, msg.CandidateTerm, nodeLogLength, msg.LogLength)
		glog.Info("Vote declined")
		return &pb.RequestForVoteReply{Ack:false}, nil
	}
}

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
	if err := s.preLoadCommitIndex(); err != nil {
		return err
	}

	glog.Infof("Starting server at %s:%d", s.nodeAddr, s.nodePort)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.nodeAddr, s.nodePort))
	if err != nil {return(err)}

	grpcServer := grpc.NewServer()
	pb.RegisterVotingServiceServer(grpcServer, &Server{})
	grpcServer.Serve(lis)
}