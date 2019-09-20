package Raft

import (
	"fmt"
)

/*
State definitions
0 - Follower
1 - Candidate
2 - Leader
*/
const (
	timeout = 1000 // ms
	state 	= 0;
)

func requestLeader() (string, int, error) {
	/*
		Entry point for RAFT. Nodes will start as slaves until timeout before moving to candidate state.
	*/
}

func Init() error {
	/*
		Entry point for RAFT. Nodes will start as slaves until timeout before moving to candidate state.
	*/
}
