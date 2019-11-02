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
	NrSuccessors	uint32			// Number of entries in finger table
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

	for baseHashGroup, successor := range c.FingerTable.Successors {
		if baseHashGroupSearched >=  baseHashGroup {
			diff = baseHashGroupSearched - baseHashGroup
		} else {
			// Correct for circular shape of hashing
			diff = baseHashGroupSearched + (c.HighestHash - baseHashGroup)
		}

		if diff == 0 {
			// Hash found
			closestHash = baseHashGroup
			glog.Infof(fmt.Sprint("Found hashgroup ", closestHash, " in table"))
			return successor
		} else if diff < smallestDiff {
			// New closest hash found
			closestHash = baseHashGroup
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
	glog.Infof(fmt.Sprint("Searching for ", baseHashGroupSearched, " from ",
						  c.BaseHashGroup))
	var closestSuccessor util.NodeInfo
	// Checks its own hash
	if baseHashGroupSearched == c.BaseHashGroup {
		glog.Infof("Found hashgroup here")
		nodeObj := util.NodeInfo{Addr: c.NodeAddr, Port: c.NodePort,
										 IsLocal: true}
		closestSuccessor = nodeObj
		util.ChordRoutingChannel.RespCh <- closestSuccessor
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
			closestSuccessor := c.search(baseHashGroupSearched)
			go func() {
				util.ChordRoutingChannel.RespCh <-closestSuccessor
			}()
			glog.Info(fmt.Sprint("Closest server ",
								 <-util.ChordRoutingChannel.ReqCh))
			glog.Info(fmt.Sprint(closestSuccessor))
		default:
		}
	}
}

func (c ChordManager) Start() {
	c.FingerTable, _ = NewFingerTable(c.NodeAddr, c.NodePort, c.NrSuccessors,
									  c.BaseHashGroup, c.HighestHash)
	c.FingerTable.FillTable()
	c.FingerTable.UpdateNodes()
	c.Routing()
}
