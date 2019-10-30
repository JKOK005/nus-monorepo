package Client

import (
	"context"
	"github.com/golang/glog"
	pb "group-project/Protobuf/Generate"
	"group-project/Utils"
)

func (c *Client) hashKey(key string) uint32 {
	/*
	Applies a hash function over a key
	Returns the hashed value of the key to represent the hash group it should be assigned to
	*/
}

func (c *Client) validate(key string) (bool, string) {
	/*
	Validates the hash group the key should be assigned to via CHORD
	We first hash the raw key to obtain a hash group number.
	We then query the Chord manager to determine if the key should be present locally or if it should be dispatched to a different address

	TODO: Implement return type of chord manager once @Johnfietsen has merged his impelmentation of chord to master
	TODO: Return types (bool, string) -> (is local node, address to be sent to if local node = False). These return types are merely a proxy and can change
	*/
	return nil, nil
}

func (c *Client) PutKey(ctx context.Context, msg *pb.PutKeyMsg) (*pb.PutKeyResp, error) {
	/*
	TODO: Check that the request is a legit request for the node, else we forward it to somewhere else
	TODO: Implement replication to slave nodes
	*/
	glog.Info("Received request to PUT key")
	Utils.PutKeyChannel.ReqCh <- msg
	isSuccess := <-Utils.PutKeyChannel.RespCh
	return &pb.PutKeyResp{Ack:isSuccess}, nil
}

func (c *Client) GetKey(ctx context.Context, msg *pb.GetKeyMsg) (*pb.GetKeyResp, error) {
	/*
	TODO: Check that the request is a legit request for the node, else we forward it to somewhere else
	TODO: Implement get key request if slave, else forward key to slave if leader
	*/
	glog.Info("Received request to GET key")
	Utils.GetKeyChannel.ReqCh <- msg.Key
	resp := <-Utils.GetKeyChannel.RespCh
	return resp, nil
}