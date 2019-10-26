package Chord

import (
	"encoding/json"
	"fmt"
	"group-project/Utils"
	"github.com/golang/glog"
)

type Protocol struct {
	MyInfo			*NodeInfo
}

func NewProtocol(myAddr string, myPort uint32, baseHashGroup uint32) (*Protocol,
																	  error) {
	if zookeeperCli, err := Utils.NewZkClient(); err != nil {
  		glog.Error(err)
  		return nil, err
  	} else {
  		glog.Infof("Setting up protocol %s:%d", myAddr, myPort)
  		nodeObj := &NodeInfo{Addr: myAddr, Port: myPort,
  							 BaseHashGroup: baseHashGroup}
  		data, _ := json.Marshal(nodeObj)
  		glog.Info(string(data))
  		err = zookeeperCli.RegisterEphemeralNode(zookeeperCli.
  					PrependNodePath(fmt.Sprintf("%d/", baseHashGroup)), data)
  		if err != nil {
  			return nil, err
  		}
  		zkCli = zookeeperCli // Cache client
  		return &Protocol{MyInfo: nodeObj}, nil
  	}
}
