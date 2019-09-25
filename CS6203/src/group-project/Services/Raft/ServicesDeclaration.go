package Raft

import (
	"context"
	"github.com/golang/glog"
	pb "group-project/Protobuf/Generate"
	util "group-project/Utils"
)

func (s *Server) requestTermNo() uint32 {
	util.GetTermNoCh.ReqCh <- true
	return <-util.GetTermNoCh.RespCh
}

func (s *Server) setTermNo(newTermNo uint32) bool {
	util.SetTermNoCh.ReqCh <- newTermNo
	return <-util.SetTermNoCh.RespCh
}

func (s *Server) RequestVote(ctx context.Context, msg *pb.RequestForVoteMsg) (*pb.RequestForVoteReply, error) {
	/*
	Evaluates if node should vote positive for a vote request RPC.
	Vote is YES if the incoming vote message V contains:
		- node_candidateTerm <= V.candidateTerm
	Vote is NO if otherwise
	*/
	termNo := s.requestTermNo()
	if termNo <= msg.CandidateTerm {
		glog.Infof("(Vote YES) %s:%d - node term no: %d <= %d", s.NodeAddr, s.NodePort, termNo, msg.CandidateTerm)
		s.setTermNo(msg.CandidateTerm)
		return &pb.RequestForVoteReply{Ack:true}, nil
	} else {
		glog.Infof("(Vote NO) %s:%d - node term no: %d > %d", s.NodeAddr, s.NodePort, termNo, msg.CandidateTerm)
		return &pb.RequestForVoteReply{Ack:false}, nil
	}
}

func (s *Server) HeartbeatCheck(ctx context.Context, msg pb.HeartBeatMsg) (*pb.HeartBeatResp, error) {
	/*
	The leader is responsible for sending heartbeats to slaves before a timeout period
	If the slave timeouts, it will promote itself to a Candidate for election
	*/
	// TODO: Renew hearbeat counter here
	return &pb.HeartBeatResp{Ack:true}, nil
}

func (s *Server) ReceiveReplication(ctx context.Context, msg pb.StatementReplicationMsg) (*pb.StatementReplicationResp, error) {
	/*
	Executes a request to set a (key, val) pair into DB
	*/
	if err := s.DbCli.Put(msg.Key, msg.Val); err != nil {
		glog.Fatal(err)
		return &pb.StatementReplicationResp{Ack:false}, err
	}
	return &pb.StatementReplicationResp{Ack:true}, nil
}

func (s *Server) ReceiveFingerTable(ctx context.Context, msg pb.FingerTableReplicationMsg) (*pb.FingerTableReplicationResp, error) {
	/*
	Executes a request to update Finger Table
	*/
	return nil, nil
}