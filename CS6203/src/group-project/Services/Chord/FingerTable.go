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
	"group-project/Services/Raft"
)

// type NodeInfo struct {
// 	Addr			string
// 	Port			uint32
// 	BaseHashGroup	uint32
// }

NodeInfo := raft.NodeInfo

type FingerTable struct {
	MyInfo			*NodeInfo
	NrSuccessors	uint32
	Successors		map[uint32]NodeInfo
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
							Successors: make(map[uint32]NodeInfo)}, nil
	}
}

func (f *FingerTable) FillTable() {
	/*
		Finds which baseHashGroups are populated (from Zookeeper)
		Choose successors and predecessor to add to the finger table
	*/
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

	f.chooseSuccessors(baseHashGroupsInt)
	// f.Predecessor = f.ChoosePredecessor(baseHashGroupsInt)
}


func (f *FingerTable) findSuccessor(baseHashGroupsInt []uint32, value uint32,
									successors map[uint32]NodeInfo) bool {
	/*
		Iterate through list of baseHashGroups
		Check if the baseHashGroup ...
	*/
	for _, eInt := range baseHashGroupsInt {
		if eInt == value {
			nodePaths, _ := zkCli.GetNodePaths(zkCli.
									PrependNodePath(fmt.Sprintf("%d", eInt)))
			// TODO determine which one is the leader
			nodeData, _ := zkCli.GetNodeValue(zkCli.PrependNodePath(fmt.
										Sprintf("%d/%s", eInt, nodePaths[0])))
			nodeInfo := new(NodeInfo)
			json.Unmarshal(nodeData, nodeInfo)
			successors[eInt] = *nodeInfo
			// successors = append(successors, *nodeInfo)
			return true
		}
	}
	return false
}


func (f *FingerTable) chooseSuccessors(baseHashGroupsInt []uint32) {

	highestBaseHashGroup := uint32(10) // will be based on hash function
	for i := uint32(0); i < f.NrSuccessors; i++ {
		value := f.MyInfo.BaseHashGroup + uint32(math.Pow(2, float64(i)))
		found := false
		for found == false {
			if value > highestBaseHashGroup {
				value = value % highestBaseHashGroup
			}
			found = f.findSuccessor(baseHashGroupsInt, value, f.Successors)
			value = value + 1
		}
	}

	glog.Info(fmt.Sprint("Successors ", f.Successors, " of ",
						 f.MyInfo.BaseHashGroup))
}

// func (f * FingerTable) ChoosePredecessor(baseHashGroupsInt []uint32) []uint32 {
//
// 	var predecessors []uint32
//
// 	p := f.MyInfo.BaseHashGroup
// 	highestBaseHashGroup := baseHashGroupsInt[len(baseHashGroupsInt)-1]
// 	found := false
// 	for found == false {
// 		if p < 0 {
// 			p = highestBaseHashGroup
// 		}
// 		for _, eInt := range baseHashGroupsInt {
// 			if eInt == uint32(p) {
// 				predecessors = append(predecessors, eInt)
// 				found = true
// 				break
// 			}
// 		}
// 		p = p - 1
// 	}
//
// 	return predecessors
// }
