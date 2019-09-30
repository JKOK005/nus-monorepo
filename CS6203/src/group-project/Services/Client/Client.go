package Client

import (
	"context"
	"github.com/golang/glog"
	pb "group-project/Protobuf/Generate"
	"group-project/Utils"
)

type Client struct {}

func (c *Client) PutKeyService(ctx context.Context, msg *pb.PutKeyMsg) (*pb.PutKeyResp, error) {
	/*
	TODO: Check that the request is a legit request for the node, else we forward it to somewhere else
	TODO: Implement replication to slave nodes
	*/
	glog.Info("Received request to PUT key")
	Utils.PutKeyChannel.ReqCh <- msg
	isSuccess := <-Utils.PutKeyChannel.RespCh
	return &pb.PutKeyResp{Ack:isSuccess}, nil
}

func (c *Client) GetKeyService(ctx context.Context, msg *pb.GetKeyMsg) (*pb.GetKeyResp, error) {
	/*
	TODO: Check that the request is a legit request for the node, else we forward it to somewhere else
	TODO: Implement get key request if slave, else forward key to slave if leader
	*/
	glog.Info("Received request to GET key")
	Utils.GetKeyChannel.ReqCh <- msg.Key
	resp := <-Utils.GetKeyChannel.RespCh
	return resp, nil
}

func (c *Client) Start() {
	/*
	TODO: Register client services
	*/
	glog.Info("Client services starting up")
}