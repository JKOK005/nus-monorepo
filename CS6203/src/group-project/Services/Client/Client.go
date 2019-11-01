package Client

import (
	"context"
	"github.com/golang/glog"
	pb "group-project/Protobuf/Generate"
	"group-project/Utils"
)

func (c *Client) locate(key string) Utils.ChannelsNodeInfo {
	/*
		Determines the hash group the key should be located at via CHORD
		We first hash the raw key to obtain a hash group number.
		We then query the Chord manager to determine if the key should be present locally or if it should be dispatched to a different address
	*/
	glog.Info("Locating key: ", key)
	hashFnct := Utils.GetHashFunction()
	Utils.ChordRoutingChannel.ReqCh <- hashFnct(key)
	return <-Utils.ChordRoutingChannel.RespCh
}

func (c *Client) forwardPut(msg *pb.PutKeyMsg, recipient Utils.ChannelsNodeInfo) (*pb.PutKeyResp, error) {
	// TODO: Placeholder for forwarding PUT requests to recipient node
	return nil, nil
}

func (c *Client) forwardGet(msg *pb.GetKeyMsg, recipient Utils.ChannelsNodeInfo) (*pb.GetKeyResp, error) {
	// TODO: Placeholder for forwarding GET requests to recipient node
	return nil, nil
}

func (c *Client) PutKey(ctx context.Context, msg *pb.PutKeyMsg) (*pb.PutKeyResp, error) {
	var isSuccess bool
	glog.Info("Received request to PUT key")
	_ = c.locate(msg.Key)
	if false {
		// TODO: Block to route request to other nodes if CHORD tells us a new address, once @Johnfiesten makes a new PR
		isSuccess = true
	} else {
		glog.Info("Received request to PUT key")
		Utils.PutKeyChannel.ReqCh <- msg
		Utils.ReplicationChannel.ReqCh <- msg
		isSuccess = <-Utils.PutKeyChannel.RespCh && <-Utils.ReplicationChannel.RespCh
	}
	return &pb.PutKeyResp{Ack:isSuccess}, nil
}

func (c *Client) GetKey(ctx context.Context, msg *pb.GetKeyMsg) (*pb.GetKeyResp, error) {
	var resp *pb.GetKeyResp
	glog.Info("Received request to GET key")
	_ = c.locate(msg.Key)
	if false {
		// TODO: Block to route request to other nodes if CHORD tells us a new address, once @Johnfiesten makes a new PR
		resp = &pb.GetKeyResp{Ack:true, Val:nil}
	} else {
		Utils.GetKeyChannel.ReqCh <- msg.Key
		resp = <-Utils.GetKeyChannel.RespCh
	}
	return resp, nil
}