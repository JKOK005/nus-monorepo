package Raft

import (
	"context"
	"encoding/json"
	"fmt"
	pb "group-project/Protobuf/Generate"
	"group-project/Utils"
	"time"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

type NodeInfo struct {
	Addr string
	Port uint32
}

type Coordinator struct {
	MyInfo            *NodeInfo
	NodesInGroup      []*NodeInfo
	PollTimeOutMs     uint32
	RefreshTimeMs     int64
	LastSyncTimeEpoch int64
}

var (
	refreshTimeMs = int64(10000)
	pollTimeOutMs = uint32(5000)
	zkCli         *Utils.SdClient
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
		nodeObj := &NodeInfo{Addr: myAddr, Port: myPort}
		data, _ := json.Marshal(nodeObj)
		glog.Info(string(data))
		err = zookeeperCli.RegisterEphemeralNode(zookeeperCli.PrependNodePath(fmt.Sprintf("%d/", baseHashGroup)), data)
		if err != nil {
			return nil, err
		}
		zkCli = zookeeperCli // Cache client
		return &Coordinator{MyInfo: nodeObj, RefreshTimeMs: refreshTimeMs, PollTimeOutMs: pollTimeOutMs}, nil
	}
}

func (c *Coordinator) marshalOne(data []byte) (*NodeInfo, error) {
	node := new(NodeInfo)
	err := json.Unmarshal(data, node)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (c *Coordinator) GetNodeZk(nodePath string) (*NodeInfo, error) {
	/*
		Returns the node value of a node registered under nodePath
	*/
	if unmarshalledNode, err := zkCli.GetNodeValue(nodePath); err != nil {
		glog.Warning(err)
		return nil, err
	} else {
		marshalledNode, _ := c.marshalOne(unmarshalledNode)
		return marshalledNode, nil
	}
}

func (c *Coordinator) RefreshNodeList(baseHashGroup uint32) error {
	// Hits Zookeeper under /<base>/baseHashGroup to fetch all nodes registered in the directory
	glog.Infof("Forced refresh of nodes in group %d requests to Zookeeper", baseHashGroup)
	if childPaths, err := zkCli.GetNodePaths(zkCli.PrependNodePath(fmt.Sprintf("%d", baseHashGroup))); err != nil {
		glog.Warning(err)
		return err
	} else {
		var nodePtrs []*NodeInfo
		for _, nodePath := range childPaths {
			if nodePtr, err := c.GetNodeZk(zkCli.PrependNodePath(fmt.Sprintf("%d/%s", baseHashGroup, nodePath))); err != nil {
				glog.Warning(err)
				return err
			} else {
				nodePtrs = append(nodePtrs, nodePtr)
			}
		}
		c.NodesInGroup = nodePtrs // Updates cache
	}
	return nil
}

func (c *Coordinator) GetNodes(baseHashGroup uint32) ([]*NodeInfo, error) {
	/*
		Returns the addresses of all nodes in the same node group

		If lastSyncTimeEpoch + refreshTimeMs < presentTimeEpoch
		- We hit Zookeeper under /<base>/baseHashGroup to fetch all nodes registered in the directory

		Else
		- Return list of nodes cached
	*/
	currTimeMs := time.Now().UnixNano() / int64(time.Millisecond)
	glog.Info(currTimeMs, " ", c.LastSyncTimeEpoch, " ", c.RefreshTimeMs)
	if currTimeMs >= c.LastSyncTimeEpoch+c.RefreshTimeMs {
		if err := c.RefreshNodeList(baseHashGroup); err != nil {
			glog.Error(err)
			return nil, err
		}
		c.LastSyncTimeEpoch = currTimeMs // Updates timestamp to ensure we do not hit ZK so soon
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.PollTimeOutMs)*time.Millisecond)
		defer conn.Close()
		defer cancel()

		client := pb.NewVotingServiceClient(conn)
		resp, err := client.RequestVote(ctx, &pb.RequestForVoteMsg{CandidateTerm: termNo})

		if ctx.Err() == context.DeadlineExceeded {
			// Request timed out. Report as timeout.
			glog.Warning("Request timed out: ", ctx.Err())
			respCh <- false
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
	respCh := make(chan bool)
	defer close(respCh)
	for _, nodePtr := range nodes {
		go c.RequestVote(nodePtr, termNo, respCh)
	}
	for range nodes {
		voteResp = append(voteResp, <-respCh)
	}
	return voteResp, nil
}

func (c *Coordinator) IssueHeartbeat(node *NodeInfo, termNo uint32, respCh chan bool) {
	/*
		Issues a heartbeat message to a node
	*/
	glog.Infof("Renewing heartbeat for addr: %s, port: %d", node.Addr, node.Port)
	if conn, err := grpc.Dial(fmt.Sprintf("%s:%d", node.Addr, node.Port), grpc.WithInsecure()); err != nil {
		glog.Warning(err)
		respCh <- false
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.PollTimeOutMs)*time.Millisecond)
		defer conn.Close()
		defer cancel()

		client := pb.NewHeartbeatServiceClient(conn)
		resp, err := client.HeartbeatCheck(ctx, &pb.HeartBeatMsg{TermNo: termNo})

		if ctx.Err() == context.DeadlineExceeded {
			// Request timed out. Report as timeout.
			glog.Warning("Request timed out: ", ctx.Err())
			respCh <- false
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

func (c *Coordinator) IssueHeartbeats(nodes []*NodeInfo, termNo uint32) ([]bool, error) {
	/*
		Issues heartbeat checks to a list of nodes
	*/
	var heartbeatsResp []bool
	respCh := make(chan bool)
	defer close(respCh)
	for _, nodePtr := range nodes {
		go c.IssueHeartbeat(nodePtr, termNo, respCh)
	}
	for range nodes {
		heartbeatsResp = append(heartbeatsResp, <-respCh)
	}
	return heartbeatsResp, nil
}

func (c *Coordinator) MarkAsLeader(baseHashGroup uint32) error {
	/*
		Node declares itself as the leader of the group
		We must:
		- Register node's address:port under /nodes/<hash_no>
		- Create a /follower/<hash_no> path if it does not exist
		- Inform all nodes registered under /follower/<hash_no> that it is the new leader
	*/
	var err error
	data, _ := json.Marshal(c.MyInfo)
	glog.Info("Register as leader: ", string(data))
	err = zkCli.SetNodeValue(zkCli.PrependNodePath(fmt.Sprintf("%d", baseHashGroup)), data)
	err = zkCli.ConstructNodesInPath(zkCli.PrependFollowerPath(fmt.Sprintf("%d", baseHashGroup)), "/", nil)
	// TODO: Inform all nodes registered under /follower/<hash_no> that it is the new leader
	return err
}
