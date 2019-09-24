package Raft

import (
	"fmt"
	"github.com/golang/glog"
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
	nodesInGroup 		[]*NodeInfo
	lastSyncTimeEpoch 	uint32
}

var (
	zkCli *Utils.SdClient
)

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

		zkCli = zookeeperCli		// Cache client
		return &Coordinator{myInfo: nodeObj}, nil
	}
}

func (c *Coordinator) marshalOne(data []byte)(*NodeInfo, error) {
	node := new(NodeInfo)
	err := json.Unmarshal(data, node)
	if err != nil {return nil, err}
	return node, nil
}

func (c *Coordinator) GetNodeZk(nodePath string) (*NodeInfo, error) {
	/*
	Returns the node value of a node registered under nodePath
	*/
	if unmarshalledNode, err := zkCli.GetNodeValue(zkCli.PrependNodePath(nodePath)); err != nil {
		glog.Fatal(err)
		return nil, err
	} else {
		marshalledNode, _ := c.marshalOne(unmarshalledNode)
		return marshalledNode, nil
	}
}

func (c *Coordinator) GetNodes(baseHashGroup uint32, refreshTimeMs uint32) ([]*NodeInfo, error) {
	/*
	Returns the addresses of all nodes in the same node group

	If lastSyncTimeEpoch + refreshTimeMs < presentTimeEpoch
	- We hit Zookeeper under /<base>/baseHashGroup to fetch all nodes registered in the directory

	Else
	- Return list of nodes cached
	*/
	currTimeMs := time.Now().Nanosecond() / 1000000

	if uint32(currTimeMs) >= c.lastSyncTimeEpoch + refreshTimeMs {
		// Hit ZK to refresh cache
		if childPaths, err := zkCli.GetNodePaths(zkCli.PrependNodePath(fmt.Sprint(baseHashGroup))); err != nil {
			glog.Fatal(err)
			return nil, err
		} else {
			var nodePtrs []*NodeInfo
			for _, nodePath := range childPaths {
				if nodePtr, err := c.GetNodeZk(zkCli.PrependNodePath(fmt.Sprint("%d/%s", baseHashGroup, nodePath))); err != nil {
					glog.Fatal(err)
					return nil, err
				} else {nodePtrs = append(nodePtrs, nodePtr)}
			}
			c.nodesInGroup = nodePtrs 					// Updates cache
			c.lastSyncTimeEpoch = uint32(currTimeMs) 	// Updates timestamp to ensure we do not hit ZK so soon
		}
	}
	return c.nodesInGroup, nil
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