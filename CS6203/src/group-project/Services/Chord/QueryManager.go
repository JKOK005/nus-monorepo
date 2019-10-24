package Chord

import (
	"fmt"
	"math"
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
}

func (q *QueryManager) HandleRequest(baseHashGroupSearched uint32) {

	found := false
	// for _, baseHashGroup := range q.FingerTable.Successors {
	// 	if baseHashGroupSearched == baseHashGroup {
	// 		found = true
	// 		glog.Infof("Found hashgroup")
	// 	}
	// }

	smallestDiff := math.Inf
	closestHash := math.Inf
	for _, baseHashGroup := range q.FingerTable.Successors {
		diff =  baseHashGroupSearched - baseHashGroup
		if diff == 0 {
			glog.Infof("Found hashgroup")
			found = true
		} else if diff < smallestDiff {
			closestHash = baseHashGroup
		}
	}
}
