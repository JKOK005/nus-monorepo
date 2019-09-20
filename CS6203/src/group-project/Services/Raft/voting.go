package Raft

import (
	"context"
	"github.com/golang/glog"
	pb "group-project/Protobuf/Generate"
)

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
	if s.termNo < msg.CandidateTerm && nodeLogLength <= msg.LogLength {
		glog.Infof("%s:%d - node term no: %d < %d & logLength: %d <= %d",
			s.nodeAddr, s.nodePort, s.termNo, msg.CandidateTerm, nodeLogLength, msg.LogLength)
		glog.Info("Vote acknowledged")
		s.setTermNo(msg.CandidateTerm)
		return &pb.RequestForVoteReply{Ack:true}, nil
	} else {
		glog.Infof("%s:%d - node term no: %d >= %d & logLength: %d > %d",
			s.nodeAddr, s.nodePort, s.termNo, msg.CandidateTerm, nodeLogLength, msg.LogLength)
		glog.Info("Vote declined")
		return &pb.RequestForVoteReply{Ack:false}, nil
	}
}