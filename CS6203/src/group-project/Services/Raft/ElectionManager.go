package Raft

import (
	"github.com/golang/glog"
	util "group-project/Utils"
	"math"
	"time"
)

type candidateState uint8

type ElectionManager struct {
	NodeAddr 		string 					// Node address
	NodePort 		uint32 					// Node port
	BaseHashGroup 	uint32 					// Base hash number for a group group (This will be registered in ZK)
	TermNo 			uint32					// Present term number
	CycleNo 		uint32					// Present cycle number
	CyclesToTimeout uint32					// We declare a timeout if cyclesToTimeout > cycleNo
	CycleTimeMs 	uint32 					// Cycle time for the start loop
	State 			candidateState 			// Present state
}

const (
	Follower 	candidateState = 0
	Candidate 	candidateState = 1
	Leader 		candidateState = 2
)

func (e *ElectionManager) setCandidateState(state candidateState) {e.State = state}
func (e *ElectionManager) setCycleNo(no uint32) {e.CycleNo = no}

func (e *ElectionManager) setTermNo(no uint32) bool {
	glog.Info("Term no set to: ", no)
	e.TermNo = no % math.MaxUint32
	return true
}

func (e *ElectionManager) votedMajority(votes []bool, quorumSize int) bool {
	yesVotes := 1 // Node votes for itself by default
	for _, vote := range votes {if vote {yesVotes++}}
	return yesVotes >= (quorumSize / 2) +1
}

func (e *ElectionManager) termNoRoutine() {
	// Handles any request for term number
	for {
		select {
		case <-util.GetTermNoCh.ReqCh: glog.Info("Term no request received"); util.GetTermNoCh.RespCh <- e.TermNo
		default:
		}
	}
}

func (e ElectionManager) Start() {
	coordCli, err := NewCoordinatorCli(e.NodeAddr, e.NodePort, e.BaseHashGroup)
	if err != nil {glog.Fatal(err); panic(err)}
	go e.termNoRoutine()

	for {
		select {
		case <- time.NewTicker(time.Duration(e.CycleTimeMs) * time.Millisecond).C:
			if e.State == Follower {
				select {
				case termNo := <-util.SetTermNoCh.ReqCh:
					glog.Infof("Setting term to ", termNo)
					util.SetTermNoCh.RespCh <- e.setTermNo(termNo)
					e.setCycleNo(0)
				default:
					if e.CycleNo > e.CyclesToTimeout {
						e.setTermNo(e.TermNo +1) // Increments term no and transit to candidate status
						e.setCandidateState(Candidate)
						e.setCycleNo(0)
					}
					glog.Info("In follower state. Cycle no: ", e.CycleNo)
				}
				e.setCycleNo(e.CycleNo +1) // Increments cycle counter in follower state
			} else if e.State == Candidate {
				glog.Info("In candidate state and requesting for votes")
				nodeLst, _ 		:= coordCli.GetNodes(e.BaseHashGroup)
				votes, _ 		:= coordCli.RequestVotes(nodeLst, e.TermNo)
				glog.Infof("Quorum size: %d, Votes are: ", len(nodeLst), votes)
				if e.votedMajority(votes, len(nodeLst)) {
					glog.Info("Received majority of votes. Promoting to leader")
					_ = coordCli.MarkAsLeader(e.BaseHashGroup)
					e.setCandidateState(Leader)
				} else {
					glog.Info("Received minority of votes. Demoting to follower")
					e.setCandidateState(Follower)
				}
			} else if e.State == Leader {
				glog.Info("In leader state")
			}
		}
	}
}