package Raft

import (
	"fmt"
	"github.com/golang/glog"
	_ "github.com/golang/glog"
	"group-project/Utils"
	"time"
	"encoding/json"
)

type NodeInfo struct {
	addr string
	port uint32
}

type Coordinator struct {
	myInfo 				*NodeInfo
	nodesInGroup 		[]NodeInfo
	lastSyncTimeEpoch 	time.Duration
}

func NewCoordinatorCli(myAddr string, myPort uint32, baseHashGroup uint32) (*Coordinator, error) {
	/*
	Creates a new Coordinator struct and returns to the user
	Also initializes node path in Zookeeper
	*/
	if zookeeperCli, err := Utils.NewZkClient(); err != nil {
		glog.Error(err)
		return nil, err
	} else {
		glog.Info(fmt.Sprintf("Addr: %s, Port: %d, hash group: %d", myAddr, myPort, baseHashGroup))
		nodeObj 	:= &NodeInfo{addr:myAddr, port:myPort}
		data, _ 	:= json.Marshal(nodeObj)
		err 		= zookeeperCli.RegisterEphemeralNode(zookeeperCli.PrependNodePath(fmt.Sprint(baseHashGroup)), data)
		if err != nil {return nil, err}
		return &Coordinator{myInfo: nodeObj}, nil
	}
}

func (c *Coordinator) GetNodes(baseHashGroup uint32, refreshTimeMs time.Duration) ([]NodeInfo, error) {
	/*
	Returns the addresses of all nodes in the same node group

	If lastSyncTimeEpoch + refreshTimeMs < presentTimeEpoch
	- We hit Zookeeper under /<base>/baseHashGroup to fetch all nodes registered in the directory

	Else
	- Return list of nodes cached
	*/
	return nil, nil
}

func (c *Coordinator) RequestVote(node NodeInfo, termNo uint32) (bool, error) {
	/*
	Requests a node for voting via gRPC
	*/
	return false, nil
}

func (c *Coordinator) RequestVotes(nodes []NodeInfo, termNo uint32) ([]bool, error) {
	/*
	Requests a list of nodes for voting via gRPC
	*/
	return nil, nil
}

func (c *Coordinator) IssueHeartbeat(node NodeInfo, termNo uint32) (bool, error) {
	/*
	Issues a heartbeat message to a node
	*/
	return false, nil
}

func (c *Coordinator) IssueHeartbeats(node []NodeInfo, termNo uint32) ([]bool, error) {
	/*
	Issues heartbeat messages to a list of nodes
	*/
	return nil, nil
}