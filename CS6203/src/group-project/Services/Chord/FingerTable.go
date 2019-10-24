package Chord

// Chord implements a faster search method
// by requiring each node to keep a finger
// table containing up to m entries,
// recall that m is the number of bits in the hash key
// https://medium.com/techlog/
// chord-building-a-dht-distributed-hash-table-in-golang-67c3ce17417b

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sort"
	"math"
	"group-project/Utils"

	"github.com/golang/glog"
)

type NodeInfo struct {
	Addr			string
	Port			uint32
	BaseHashGroup	uint32
}

type FingerTable struct {
	MyInfo			*NodeInfo
	NrSuccessors	uint32
	Predecessor		[]uint32
	Successors		[]uint32
}

var (
	zkCli 			*Utils.SdClient
)

func NewFingerTable(myAddr string, myPort uint32, nrSuccessors uint32,
					baseHashGroup uint32) (*FingerTable, error) {
	/*
		Creates a new FingerTable struct and returns to the user
		Also initializes node path in Zookeeper
	*/
	if zookeeperCli, err := Utils.NewZkClient(); err != nil {
		glog.Error(err)
		return nil, err
	} else {
		glog.Infof("Building finger table %s:%d", myAddr, myPort)
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
		return &FingerTable{MyInfo: nodeObj, NrSuccessors: nrSuccessors,
							Predecessor: nil, Successors: nil}, nil
	}
}

func (f *FingerTable) FillTable() {

	baseHashGroupsFound, _ := zkCli.GetNodePaths("/nodes")
	var baseHashGroupsPopulated []string

	for _, baseHashGroup := range baseHashGroupsFound {
		nodesFound, _ := zkCli.GetNodePaths(fmt.Sprint("/nodes/",
													   baseHashGroup))
		if len(nodesFound) > 0 {
			baseHashGroupsPopulated = append(baseHashGroupsPopulated,
											 baseHashGroup)
		}
	}

	var baseHashGroupsInt []uint32
	for _, eStr := range baseHashGroupsPopulated {
		eInt, _ := strconv.ParseUint(eStr, 10, 32)
		baseHashGroupsInt = append(baseHashGroupsInt, uint32(eInt))
	}
	sort.Slice(baseHashGroupsInt, func(i,
		j int) bool { return baseHashGroupsInt[i] < baseHashGroupsInt[j] })
	glog.Info("BaseHashGroups populated: ", baseHashGroupsInt)
	f.Successors = f.ChooseSuccessors(baseHashGroupsInt)
	f.Predecessor = f.ChoosePredecessor(baseHashGroupsInt)
}

func (f *FingerTable) ChooseSuccessors(baseHashGroupsInt []uint32) []uint32 {

	var successors []uint32

	// nrSuccessors := uint32(len(baseHashGroupsInt) / 2)
	highestBaseHashGroup := baseHashGroupsInt[len(baseHashGroupsInt)-1]
	for i := uint32(0); i < f.NrSuccessors; i++ {
		p := f.MyInfo.BaseHashGroup + uint32(math.Pow(2, float64(i)))
		found := false
		for found == false {
			if p > highestBaseHashGroup {
				p = p % highestBaseHashGroup
			}
			for _, eInt := range baseHashGroupsInt {
				if eInt == uint32(p) {
					// glog.Infof("match")
					successors = append(successors, eInt)
					found = true
					break
				}
			}
			p = p + 1
		}
	}

	glog.Info("Successors chosen: ", successors)

	return successors
}

func (f * FingerTable) ChoosePredecessor(baseHashGroupsInt []uint32) []uint32 {

	var predecessors []uint32

	p := f.MyInfo.BaseHashGroup
	highestBaseHashGroup := baseHashGroupsInt[len(baseHashGroupsInt)-1]
	found := false
	for found == false {
		if p < 0 {
			p = highestBaseHashGroup
		}
		for _, eInt := range baseHashGroupsInt {
			if eInt == uint32(p) {
				// eStr := strconv.FormatUint(uint64(eInt), 10)
				predecessors = append(predecessors, eInt)
				found = true
				break
			}
		}
		p = p - 1
	}

	return predecessors
}
