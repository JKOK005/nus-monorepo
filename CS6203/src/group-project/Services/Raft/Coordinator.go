package Raft

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	pb "group-project/Protobuf/Generate"
	"group-project/Utils"
	"time"
)

type NodeInfo struct {
	Addr string
	Port uint32
}

type Coordinator struct {
	MyInfo 				*NodeInfo
	NodesInGroup 		[]*NodeInfo
	RefreshTimeMs 		uint32
	LastSyncTimeEpoch 	uint32
	PollTimeOutMs 		uint32
}

var (
	refreshTimeMs = uint32(10000)
	pollTimeOutMs = uint32(5000)
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
		nodeObj 	:= &NodeInfo{Addr:myAddr, Port:myPort}
		data, _ 	:= json.Marshal(nodeObj)
		glog.Info(string(data))
		err 		= zookeeperCli.RegisterEphemeralNode(zookeeperCli.PrependNodePath(fmt.Sprintf("%d/", baseHashGroup)), data)
		if err != nil {return nil, err}

		zkCli = zookeeperCli		// Cache client
		return &Coordinator{MyInfo: nodeObj, RefreshTimeMs: refreshTimeMs, PollTimeOutMs: pollTimeOutMs}, nil
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
	glog.Info(nodePath)
	if unmarshalledNode, err := zkCli.GetNodeValue(nodePath); err != nil {
		glog.Warning(err)
		return nil, err
	} else {
		marshalledNode, _ := c.marshalOne(unmarshalledNode)
		return marshalledNode, nil
	}
}

func (c *Coordinator) GetNodes(baseHashGroup uint32) ([]*NodeInfo, error) {
	/*
	Returns the addresses of all nodes in the same node group

	If lastSyncTimeEpoch + refreshTimeMs < presentTimeEpoch
	- We hit Zookeeper under /<base>/baseHashGroup to fetch all nodes registered in the directory

	Else
	- Return list of nodes cached
	*/
	currTimeMs := time.Now().Nanosecond() * 1000000

	if uint32(currTimeMs) >= c.LastSyncTimeEpoch + c.RefreshTimeMs {
		// Hit ZK to refresh cache
		if childPaths, err := zkCli.GetNodePaths(zkCli.PrependNodePath(fmt.Sprintf("%d", baseHashGroup))); err != nil {
			glog.Warning(err)
			return nil, err
		} else {
			var nodePtrs []*NodeInfo
			for _, nodePath := range childPaths {
				if nodePtr, err := c.GetNodeZk(zkCli.PrependNodePath(fmt.Sprintf("%d/%s", baseHashGroup, nodePath))); err != nil {
					glog.Warning(err)
					return nil, err
				} else {nodePtrs = append(nodePtrs, nodePtr)}
			}
			c.NodesInGroup = nodePtrs 					// Updates cache
			c.LastSyncTimeEpoch = uint32(currTimeMs) 	// Updates timestamp to ensure we do not hit ZK so soon
		}
	}
	return c.NodesInGroup, nil
}

func (c *Coordinator) RequestVote(node *NodeInfo, termNo uint32, respCh chan bool) {
	/*
	Requests a node for voting via gRPC
	*/
	glog.Infof("Requesting vote from addr: %s, port: %d", node.Addr, node.Port)
	if conn, err := grpc.Dial(fmt.Sprintf("%s:%d", node.Addr, node.Port), grpc.WithInsecure()); err != nil {
		glog.Warning(err)
		respCh <- false
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.PollTimeOutMs) * time.Millisecond)
		defer conn.Close(); defer cancel()

		client := pb.NewVotingServiceClient(conn)
		resp, err := client.RequestVote(ctx, &pb.RequestForVoteMsg{CandidateTerm: termNo})

		if ctx.Err() == context.DeadlineExceeded {
			// Request timed out. Report as timeout.
			glog.Warning("Request timed out: ", ctx.Err())
		} else {
			if err != nil {
				glog.Warning(err)
				respCh <- false
			} else {
				glog.Infof("Received response: %t from addr: %s, port: %d", resp.Ack, node.Addr, node.Port)
				respCh <- resp.Ack
			}
		}
	}
}

func (c *Coordinator) RequestVotes(nodes []*NodeInfo, termNo uint32) ([]bool, error) {
	/*
	Requests a list of nodes for voting via gRPC
	*/
	var voteResp []bool
	respCh := make (chan bool)
	defer close(respCh)
	for _, nodePtr := range nodes {go c.RequestVote(nodePtr, termNo, respCh)}
	for range nodes {voteResp = append(voteResp, <-respCh)}
	return voteResp, nil
}

func (c *Coordinator) IssueHeartbeat(node *NodeInfo, termNo uint32) (bool, error) {
	/*
	Issues a heartbeat message to a node
	*/
	return false, nil
}

func (c *Coordinator) IssueHeartbeats(node []*NodeInfo, termNo uint32) ([]bool, error) {
	/*
	Issues heartbeat messages to a list of nodes
	*/
	return nil, nil
}