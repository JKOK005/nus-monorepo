package Raft

import (
	"github.com/golang/glog"
	"time"
)

type NodeInfo struct {
	addr string
	port uint32
}

type Coordinator struct {
	nodes []NodeInfo
}

var (
	lastSyncTimeEpoch = time.Duration(0)
)

func (c *Coordinator) GetNodes(baseHashGroup uint32, refreshTimeMs time.Duration) ([]NodeInfo, error) {
	/*
	Returns the addresses of all nodes in the same node group

	If lastSyncTimeEpoch + refreshTimeMs < presentTimeEpoch
	- We hit Zookeeper under /<base>/baseHashGroup to fetch all nodes registered in the directory

	Else
	- Return list of nodes cached
	*/
}

func (c *Coordinator) RequestVote(node NodeInfo, termNo uint32) (bool, error) {
	/*
	Requests a node for voting via gRPC
	*/
}

func (c *Coordinator) RequestVotes(nodes []NodeInfo, termNo uint32) ([]bool, error) {
	/*
	Requests a list of nodes for voting via gRPC
	*/
}

func (c *Coordinator) IssueHeartbeat(node NodeInfo, termNo uint32) (bool, error) {
	/*
	Issues a heartbeat message to a node
	*/
}

func (c *Coordinator) IssueHeartbeats(node []NodeInfo, termNo uint32) ([]bool, error) {
	/*
	Issues heartbeat messages to a list of nodes
	*/
}