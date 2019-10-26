package Chord

import (
	util "group-project/Utils"

	"fmt"
	"github.com/golang/glog"
)

type QueryManager struct {
	NodeAddr		string
	NodePort		uint32
	BaseHashGroup	uint32
	FingerTable		*FingerTable
}

func (q QueryManager) Start() {
	q.FingerTable, _ = NewFingerTable(q.NodeAddr, q.NodePort, 3,
									  q.BaseHashGroup)
	q.FingerTable.FillTable()
	q.HandleRequests()
	q.HandlePuts()
}

func (q *QueryManager) searchFingerTable(baseHashGroupSearched uint32) {

	smallestDiff := uint32(1e9)
	closestHash := uint32(1e9)
	highest := uint32(10)
	var diff uint32

	for _, baseHashGroupFound := range q.FingerTable.Successors {
		if baseHashGroupSearched >= baseHashGroupFound {
			diff = baseHashGroupSearched - baseHashGroupFound
		} else {
			diff = baseHashGroupSearched + (highest - baseHashGroupFound)
		}

		glog.Infof(fmt.Sprint(baseHashGroupSearched, baseHashGroupFound, diff))

		if diff == 0 {
			closestHash = baseHashGroupFound
			glog.Infof(fmt.Sprint("Found hashgroup ", closestHash, " in table"))
			return
		} else if diff < smallestDiff {
			closestHash = baseHashGroupFound
			glog.Infof(fmt.Sprint("Found closer hash ", closestHash))
		}
	}
	glog.Infof(fmt.Sprint("Redirecting to ", closestHash))
}


func (q *QueryManager) search(baseHashGroupSearched uint32) {
	glog.Infof(fmt.Sprint("Searching for ", baseHashGroupSearched,
						  " from ", q.BaseHashGroup))

	if baseHashGroupSearched == q.BaseHashGroup {
		glog.Infof("Found hashgroup here")
	} else {
		q.searchFingerTable(baseHashGroupSearched)
	}
}

func (q *QueryManager) HandleRequests() {

	for {
		select {
		case baseHashGroupSearched := <-util.SetRequestChannel.ReqCh: {
			q.search(baseHashGroupSearched)
		}
		default:
		}
	}

}

func (q *QueryManager) HandlePuts() {
	for {
		select {
		case baseHashGroupSearched := <-util.SetPutChannel.ReqCh: {
			q.search(baseHashGroupSearched)
		}
		default:
		}
	}
}
