package Chord

import (
	util "group-project/Utils"

	"fmt"
	"github.com/golang/glog"
	// "google.golang.org/grpc"
)

type ChordManager struct {
	NodeAddr		string
	NodePort		uint32
	BaseHashGroup	uint32
	FingerTable		*FingerTable
}

func (c ChordManager) Start() {
	c.FingerTable, _ = NewFingerTable(c.NodeAddr, c.NodePort, 3,
									  c.BaseHashGroup)
	c.FingerTable.FillTable()
	c.HandleRequests()
	c.HandlePuts()
}

func (c *ChordManager) searchFingerTable(baseHashGroupSearched uint32) {

	smallestDiff := uint32(1e9)
	closestHash := uint32(1e9)
	highest := uint32(10)
	var diff uint32

	for _, successor := range c.FingerTable.Successors {
		if baseHashGroupSearched >= successor.BaseHashGroup {
			diff = baseHashGroupSearched - successor.BaseHashGroup
		} else {
			diff = baseHashGroupSearched + (highest - successor.BaseHashGroup)
		}

		if diff == 0 {
			closestHash = successor.BaseHashGroup
			glog.Infof(fmt.Sprint("Found hashgroup ", closestHash, " in table"))
			return
		} else if diff < smallestDiff {
			closestHash = successor.BaseHashGroup
			glog.Infof(fmt.Sprint("Found closer hash ", closestHash))
		}
	}
	glog.Infof(fmt.Sprint("Redirecting to ", closestHash))
}


func (c *ChordManager) search(baseHashGroupSearched uint32) {
	glog.Infof(fmt.Sprint("Searching for ", baseHashGroupSearched,
						  " from ", c.BaseHashGroup))

	if baseHashGroupSearched == c.BaseHashGroup {
		glog.Infof("Found hashgroup here")
	} else {
		c.searchFingerTable(baseHashGroupSearched)
	}
}

func (c *ChordManager) HandleRequests() {
	for {
		select {
		case baseHashGroupSearched := <-util.SetRequestChannel.ReqCh: {
			c.search(baseHashGroupSearched)
		}
		default:
		}
	}
}

func (c *ChordManager) HandlePuts() {
	for {
		select {
		case baseHashGroupSearched := <-util.SetPutChannel.ReqCh: {
			c.search(baseHashGroupSearched)
		}
		default:
		}
	}
}
