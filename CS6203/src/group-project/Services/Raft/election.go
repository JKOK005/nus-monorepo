package Raft

import (
	"errors"
	"github.com/golang/glog"
	"group-project/Utils"
	"math"
	"time"
)

type candiddateState uint8

type ElectionManager struct {
	baseHashGroup 	uint32 					// Base hash number for a group group (This will be registered in ZK)
	termNo 			uint32					// Present term number
	cycleNo 		uint32					// Present cycle number
	cyclesToTimeout uint32					// We declare a timeout if cyclesToTimeout > cycleNo
	cycleTimeMs 	uint32 					// Cycle time for the start loop
	state 			candiddateState 		// Present state
	termNoChannel 	*Utils.TermNoChannel
}

const (
	follower 	candiddateState = 0
	candidate 	candiddateState = 1
	leader 		candiddateState = 2
)

func (e *ElectionManager) setCandidateState(state candiddateState) {e.state = state}
func (e *ElectionManager) setCycleNo(no uint32) {e.cycleNo = no}

func (e *ElectionManager) setTermNo(no uint32) bool {
	e.termNo = no % math.MaxUint32
	return true
}

func (e *ElectionManager) Start() error {
	//loopStartTimeMs := time.Now().Nanosecond() / 1000000
	for {
		select {
		case <- time.NewTicker(time.Duration(e.cycleTimeMs) * time.Millisecond).C:
			if e.state == follower {
				select {
				case termNo := <-e.termNoChannel.ReqCh:
					glog.Infof("ElectionManager: Setting term to ", termNo)
					e.setCycleNo(0)
					e.termNoChannel.RespCh <- e.setTermNo(termNo)
				default:
					if e.cycleNo > e.cyclesToTimeout {
						e.setTermNo(e.termNo +1) // Increments term no and transit to candidate status
						e.setCandidateState(candidate)
					}
					glog.Info("ElectionManager: In follower state")
				}
				e.setCycleNo(e.cycleNo +1) // Increments cycle counter in follower state
			} else if e.state == candidate {
				glog.Info("I am in candidate")
			} else if e.state == leader {
				glog.Info("I am in leader")
			} else {
				return errors.New("ElectionManager: state is not allowed. How did the logic end up in this statement")
			}
		}
	}
	//sleepTime := int(e.cycleTimeMs) + loopStartTimeMs - (time.Now().Nanosecond() / 1000000)
	//glog.Infof("ElectionManager: Sleeping for %d ms", sleepTime)
	//time.Sleep(time.Duration(sleepTime) * time.Millisecond)
}