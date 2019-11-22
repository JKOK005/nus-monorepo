package Client

import (
	"fmt"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	pb "group-project/Protobuf/Generate"
	dep "group-project/Utils"
	"net"
)

type Client struct {
	NodeAddr 	string
	NodePort 	uint32
}

func (c Client) Start() {
	/*
		Starts the RAFT GRPC server that listens to incoming requests

		TODO: 	Do not hard code localhost into server start :)- .... but no time to fix this :'(
				Hard coding hack done as I just realized that clients have to listen to localhost but register their own address
				differently in ZK when we use other deployment methods such as Docker Compose or Kubenetes
				The difference between LISTENER_ADR and REGISTER_LISTENER_DNS environment variables is as follows:
				- REGISTER_LISTENER_DNS: Address to be registered into ZK for discovery
				- LISTENER_ADR: Address that clients subscribe to to listen to requests
	*/
	registerAddr := dep.GetEnvStr("LISTENER_ADR", c.NodeAddr)
	glog.Infof("Starting client at %s:%d", registerAddr, c.NodePort)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", registerAddr, c.NodePort))
	if err != nil {glog.Fatal(err); panic(err)}

	grpcServer := grpc.NewServer()
	pb.RegisterGetKeyServiceServer(grpcServer, &c)
	pb.RegisterPutKeyServiceServer(grpcServer, &c)
	pb.RegisterRecomputeFingerTableServiceServer(grpcServer, &c)
	grpcServer.Serve(lis)
}