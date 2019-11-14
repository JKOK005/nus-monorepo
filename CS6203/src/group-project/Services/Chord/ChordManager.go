 package Chord

import (
	util "group-project/Utils"

	"fmt"
	"github.com/golang/glog"
)

type ChordManager struct {
	NodeAddr		string			// Node address
	NodePort		uint32			// Node port
	BaseHashGroup	uint32			// Base hash number for a group
	FingerTable		*FingerTable	// Struct containing node info of successors
	HighestHash		uint32			// Highest possible hash value
}

func (c *ChordManager) searchFingerTable(baseHashGroupSearched uint32) util.NodeInfo {
	/*
		Receives a bashashgroup it needs to find
		Iterates through the finger table and looks for closest hash
		Returns the node information of the server with the closest hash
	*/
	smallestDiff := uint32(1e9)
	var closestHash uint32
	var diff uint32
	var closestSuccessor util.NodeInfo

	for _, successor := range c.FingerTable.Successors {
		if baseHashGroupSearched >=  successor.BaseHashGroup {
			diff = baseHashGroupSearched - successor.BaseHashGroup
		} else {
			// Correct for circular shape of hashing
			diff = baseHashGroupSearched + (c.HighestHash - successor.BaseHashGroup)
		}

		if diff == 0 {
			// Hash found
			closestHash = successor.BaseHashGroup
			glog.Infof(fmt.Sprint("Found hashgroup ", closestHash, " in table"))
			return successor
		} else if diff < smallestDiff {
			// New closest hash found
			closestHash = successor.BaseHashGroup
			closestSuccessor = successor
			glog.Infof(fmt.Sprint("Found closer hash ", closestHash))
		}
		closestSuccessor = successor
	}
	glog.Infof(fmt.Sprint("Redirecting to ", closestHash))

	return closestSuccessor
}


func (c *ChordManager) search(baseHashGroupSearched uint32) util.NodeInfo {
	/*
		Receives a bashashgroup it needs to find
		Checks if this server contains the hash, checks finger table otherwise
		Returns the node information of the server with the closest hash
	*/
	glog.Infof(fmt.Sprint("Searching for ", baseHashGroupSearched, " from ", c.BaseHashGroup))
	var closestSuccessor util.NodeInfo

	// Checks its own hash
	if len(c.FingerTable.Successors) == 0 ||
	   baseHashGroupSearched == c.BaseHashGroup ||
	   baseHashGroupSearched < c.FingerTable.Successors[0].BaseHashGroup {
	   	nodeObj := util.NodeInfo{Addr: c.NodeAddr, Port: c.NodePort,
								 BaseHashGroup : c.BaseHashGroup, IsLocal: true}
		closestSuccessor = nodeObj
		glog.Infof("Hashgroup belongs to this node")
	} else {
		closestSuccessor = c.searchFingerTable(baseHashGroupSearched)
	}
	return closestSuccessor
}

func (c *ChordManager) Routing() {
	/*
		Routine that constantly checks ChordRoutingChannel.ReqCh
		If there is a request, returns the nodeinfo of closest server
	*/
	for {
		select {
		case baseHashGroupSearched := <-util.ChordRoutingChannel.ReqCh:
			glog.Info(fmt.Sprint("Route for ", baseHashGroupSearched))
			closestSuccessor := c.search(baseHashGroupSearched)
			glog.Info(fmt.Sprint("Closest server: ", closestSuccessor))
			util.ChordRoutingChannel.RespCh <-closestSuccessor
		default:
		}
	}
}

func (c ChordManager) Start() {
	c.FingerTable, _ = NewFingerTable(c.NodeAddr, c.NodePort, c.BaseHashGroup,
									  c.HighestHash)
	c.FingerTable.FillTable()
	go c.FingerTable.UpdateNodes()
	go c.Routing()
}
