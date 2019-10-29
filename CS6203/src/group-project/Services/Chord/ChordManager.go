package Chord

import (
	util "group-project/Utils"

	"fmt"
	"github.com/golang/glog"
	// "google.golang.org/grpc"
	// "group-project/Services/Election"
)

type ChordManager struct {
	NodeAddr		string
	NodePort		uint32
	BaseHashGroup	uint32
	NrSuccessors	uint32
	FingerTable		*FingerTable
}

func (c ChordManager) Start() {
	c.FingerTable, _ = NewFingerTable(c.NodeAddr, c.NodePort, c.NrSuccessors,
									  c.BaseHashGroup)
	c.FingerTable.FillTable()
	c.Routing()
	glog.Info("sdfsdfsdfsdfsdf")
}

func (c *ChordManager) searchFingerTable(baseHashGroupSearched uint32) util.ChannelsNodeInfo {

	smallestDiff := uint32(1e9)
	closestHash := uint32(1e9)
	highest := uint32(10)
	var diff uint32
	var closestSuccessor util.ChannelsNodeInfo

	for baseHashGroup, successor := range c.FingerTable.Successors {
		if baseHashGroupSearched >=  baseHashGroup {
			diff = baseHashGroupSearched - baseHashGroup
		} else {
			diff = baseHashGroupSearched + (highest - baseHashGroup)
		}

		if diff == 0 {
			closestHash = baseHashGroup
			glog.Infof(fmt.Sprint("Found hashgroup ", closestHash, " in table"))
			return successor
		} else if diff < smallestDiff {
			closestHash = baseHashGroup
			closestSuccessor = successor
			glog.Infof(fmt.Sprint("Found closer hash ", successor))
		}
		closestSuccessor = successor
	}
	glog.Infof(fmt.Sprint("Redirecting to ", closestHash))

	return closestSuccessor
}


func (c *ChordManager) search(baseHashGroupSearched uint32) util.ChannelsNodeInfo {
	glog.Infof(fmt.Sprint("Searching for ", baseHashGroupSearched,
						  " from ", c.BaseHashGroup))

	var closestSuccessor util.ChannelsNodeInfo

	if baseHashGroupSearched == c.BaseHashGroup {
		glog.Infof("Found hashgroup here")
		nodeObj := util.ChannelsNodeInfo{Addr: c.NodeAddr, Port: c.NodePort}
		closestSuccessor = nodeObj
	} else {
		closestSuccessor = c.searchFingerTable(baseHashGroupSearched)
	}
	return closestSuccessor
}

func (c *ChordManager) Routing() {
	for {
		select {
		case baseHashGroupSearched := <-util.ChordRoutingChannel.RespCh: {
			closestSuccessor := c.search(2)
			glog.Info("\n\n\n\n 1", baseHashGroupSearched)
			glog.Info(closestSuccessor)
			util.ChordRoutingChannel.ReqCh <- 2
			glog.Info("\n\n\n\n 2")
			glog.Info(fmt.Sprint("Closest server ",
								 <-util.ChordRoutingChannel.ReqCh))
			glog.Info(fmt.Sprint(closestSuccessor))
		}
		default:
		}
	}
}
