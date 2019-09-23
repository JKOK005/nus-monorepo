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

		Vote is NO if otherwise
	*/
	if s.termNo <= msg.CandidateTerm {
		glog.Infof("(Vote YES) %s:%d - node term no: %d <= %d", s.nodeAddr, s.nodePort, s.termNo, msg.CandidateTerm)
		s.setTermNo(msg.CandidateTerm)
		return &pb.RequestForVoteReply{Ack:true}, nil
	} else {
		glog.Infof("(Vote NO) %s:%d - node term no: %d > %d", s.nodeAddr, s.nodePort, s.termNo, msg.CandidateTerm)
		return &pb.RequestForVoteReply{Ack:false}, nil
	}
}