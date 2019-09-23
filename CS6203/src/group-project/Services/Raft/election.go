package Raft

import (
	"github.com/golang/glog"
	"group-project/Utils"
	"math"
	"time"
)

type candiddateState uint8

type ElectionManager struct {
	BaseHashGroup 	uint32 					// Base hash number for a group group (This will be registered in ZK)
	TermNo 			uint32					// Present term number
	CycleNo 		uint32					// Present cycle number
	CyclesToTimeout uint32					// We declare a timeout if cyclesToTimeout > cycleNo
	CycleTimeMs 	uint32 					// Cycle time for the start loop
	State 			candiddateState 		// Present state
	TermNoChannel 	*Utils.TermNoChannel
}

const (
	Follower 	candiddateState = 0
	Candidate 	candiddateState = 1
	Leader 		candiddateState = 2
)

func (e *ElectionManager) setCandidateState(state candiddateState) {e.State = state}
func (e *ElectionManager) setCycleNo(no uint32) {e.CycleNo = no}

func (e *ElectionManager) setTermNo(no uint32) bool {
	e.TermNo = no % math.MaxUint32
	return true
}

func (e *ElectionManager) Start() {
	//loopStartTimeMs := time.Now().Nanosecond() / 1000000
	for {
		select {
		case <- time.NewTicker(time.Duration(e.CycleTimeMs) * time.Millisecond).C:
			if e.State == Follower {
				select {
				case termNo := <-e.TermNoChannel.ReqCh:
					glog.Infof("Setting term to ", termNo)
					e.setCycleNo(0)
					e.TermNoChannel.RespCh <- e.setTermNo(termNo)
				default:
					if e.CycleNo > e.CyclesToTimeout {
						e.setTermNo(e.TermNo +1) // Increments term no and transit to candidate status
						e.setCandidateState(Candidate)
					}
					glog.Info("In follower state")
				}
				e.setCycleNo(e.CycleNo +1) // Increments cycle counter in follower state
			} else if e.State == Candidate {
				glog.Info("I am in candidate")
			} else if e.State == Leader {
				glog.Info("I am in leader")
			}
		}
	}
	//sleepTime := int(e.cycleTimeMs) + loopStartTimeMs - (time.Now().Nanosecond() / 1000000)
	//glog.Infof("ElectionManager: Sleeping for %d ms", sleepTime)
	//time.Sleep(time.Duration(sleepTime) * time.Millisecond)
}