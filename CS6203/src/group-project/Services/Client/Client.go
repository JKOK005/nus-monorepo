package Client

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	pb "group-project/Protobuf/Generate"
	"group-project/Utils"
	"time"
)

func (c *Client) locate(key string) Utils.NodeInfo {
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

var (
	pollTimeOutMs = uint32(5000)
)

func (c *Client) forwardPut(msg *pb.PutKeyMsg, recipient Utils.NodeInfo) (*pb.PutKeyResp, error) {
	/*
		Forwards a PUT request to the nodes's client port

		TODO: 	Remove hack for (node.Port +1), since client ports are +1 from port by default
					Got to think of a better way to communicate with the client moving forward
	*/
	fmt.Println("Forwarding put: ", recipient)
	if conn, err := grpc.Dial(fmt.Sprintf("%s:%d", recipient.Addr, recipient.Port +1), grpc.WithInsecure()); err != nil {
		glog.Error(err)
		return nil, err
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pollTimeOutMs) * time.Millisecond)
		defer conn.Close(); defer cancel()

		client := pb.NewPutKeyServiceClient(conn)
		if resp, err := client.PutKey(ctx, msg); err != nil {
			glog.Error(err)
			return nil, err
		} else {
			resp.Stats.NoOfHops += 1
			return resp, nil
		}
	}
}

func (c *Client) forwardGet(msg *pb.GetKeyMsg, recipient Utils.NodeInfo) (*pb.GetKeyResp, error) {
	/*
		Forwards a GET request to the nodes's client port

		TODO: 	Remove hack for (node.Port +1), since client ports are +1 from port by default
					Got to think of a better way to communicate with the client moving forward
	*/
	fmt.Println("Forwarding get: ", recipient)
	if conn, err := grpc.Dial(fmt.Sprintf("%s:%d", recipient.Addr, recipient.Port +1), grpc.WithInsecure()); err != nil {
		glog.Error(err)
		return nil, err
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pollTimeOutMs) * time.Millisecond)
		defer conn.Close(); defer cancel()

		client := pb.NewGetKeyServiceClient(conn)
		if resp, err := client.GetKey(ctx, msg); err != nil {
			glog.Error(err)
			return nil, err
		} else {
			resp.Stats.NoOfHops += 1
			return resp, nil
		}
	}
}

func (c *Client) PutKey(ctx context.Context, msg *pb.PutKeyMsg) (*pb.PutKeyResp, error) {
	var resp *pb.PutKeyResp
	glog.Info("Received request to PUT key")
	routeToNode := c.locate(msg.Key)
	if !routeToNode.IsLocal {
		if attempt, err := c.forwardPut(msg, routeToNode); err != nil {
			glog.Error(err)
			resp = &pb.PutKeyResp{}
		} else{
			resp = attempt
		}
	} else {
		fmt.Println("Put key ", msg)
		Utils.PutKeyChannel.ReqCh <- msg
		Utils.ReplicationChannel.ReqCh <- msg
		resp = &pb.PutKeyResp{	Ack: <-Utils.PutKeyChannel.RespCh && <-Utils.ReplicationChannel.RespCh,
								Stats: &pb.MsgStats{NoOfHops:0}}
	}
	return resp, nil
}

func (c *Client) GetKey(ctx context.Context, msg *pb.GetKeyMsg) (*pb.GetKeyResp, error) {
	var resp *pb.GetKeyResp
	glog.Info("Received request to GET key")
	routeToNode := c.locate(msg.Key)
	if !routeToNode.IsLocal {
		if attempt, err := c.forwardGet(msg, routeToNode); err != nil {
			glog.Error(err)
			resp = &pb.GetKeyResp{}
		} else{
			resp = attempt
		}
	} else {
		fmt.Println("Get key ", msg)
		Utils.GetKeyChannel.ReqCh <- msg.Key
		resp = <-Utils.GetKeyChannel.RespCh
		resp.Stats = &pb.MsgStats{NoOfHops: 0}
	}
	return resp, nil
}

func (c *Client) RecomputeFingerTable(ctx context.Context, msg *pb.RecomputeFingerTableMsg) (*pb.RecomputeFingerTableResp, error){
	/*
		Exposes service to recompute the finger table

		This is used by external leaders who want to instruct others that something has changed in the hash ring
		, hence they need to recompute their finger table to remain up to date
	*/
	glog.Info("Received request to recompute finger table")
	Utils.ChordUpdateChannel.ReqCh <- true
	return &pb.RecomputeFingerTableResp{Ack: <-Utils.ChordUpdateChannel.RespCh}, nil
}