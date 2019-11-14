package Chord

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sort"
	"math"
	util "group-project/Utils"
	"github.com/golang/glog"
)

type FingerTable struct {
	MyInfo			*util.NodeInfo		// Node info
	NrSuccessors	uint32				// Number of entries in finger table
	Successors		[]util.NodeInfo	// Entries in finger table
	HighestHash		uint32				// Highest possible hash value
}

var (
	zkCli 			*util.SdClient
)


func NewFingerTable(myAddr string, myPort uint32, baseHashGroup uint32,
					highestHash uint32) (*FingerTable, error) {
	/*
		Creates a new FingerTable struct and returns to the user
		Also initializes node path in Zookeeper
	*/
	if zookeeperCli, err := util.NewZkClient(); err != nil {
		glog.Error(err)
		return nil, err
	} else {
		glog.Infof("Building finger table %s:%d", myAddr, myPort)
		nodeObj := &util.NodeInfo{Addr: myAddr, Port: myPort,
								  BaseHashGroup : baseHashGroup, IsLocal: true}
		data, _ := json.Marshal(nodeObj)
		glog.Info(string(data))
		err = zookeeperCli.RegisterEphemeralNode(zookeeperCli.
					PrependFollowerPath(fmt.Sprintf("%d/", baseHashGroup)), data)
		if err != nil {
			return nil, err
		}
		zkCli = zookeeperCli // Cache client
		var emptySuccessors []util.NodeInfo
		return &FingerTable{MyInfo: nodeObj, Successors: emptySuccessors,
							HighestHash: highestHash}, nil
	}
}

func (f *FingerTable) findSuccessor(baseHashGroupsInt []uint32, value uint32,
									successors *[]util.NodeInfo) bool {
	/*
		Iterate through list of baseHashGroups
		Check if the correct hash is found and add to the successors
	*/
	for _, eInt := range baseHashGroupsInt {
		if eInt == value {
			nodePath, _ := zkCli.GetNodePaths(zkCli.
									PrependNodePath(fmt.Sprintf("%d", eInt)))
			nodeData, _ := zkCli.GetNodeValue(zkCli.PrependNodePath(fmt.
										Sprintf("%d/%s", eInt, nodePath[0])))
			nodeInfo := new(util.NodeInfo)
			json.Unmarshal(nodeData, nodeInfo)
			nodeInfo.BaseHashGroup = eInt
			if nodeInfo.Port == f.MyInfo.Port &&
			   nodeInfo.Addr == f.MyInfo.Addr {
				nodeInfo.IsLocal = true
			} else {
				nodeInfo.IsLocal = false
			}
			*successors = append(*successors, *nodeInfo)
			return true
		}
	}
	return false
}


func (f *FingerTable) chooseSuccessors(baseHashGroupsInt []uint32) {
	/*
		Finds baseHashGroup closest to the baseHashGroup needed
		Calls findSuccessor to add it to the list
	*/
	for i := uint32(0); i < f.NrSuccessors; i++ {
		value := f.MyInfo.BaseHashGroup + uint32(math.Pow(2, float64(i)))
		found := false
		for found == false {
			if value > f.HighestHash {
				value = value % f.HighestHash
			}
			found = f.findSuccessor(baseHashGroupsInt, value, &f.Successors)
			value = value + 1
		}
	}
	glog.Info(fmt.Sprint("Successors ", f.Successors, " of ", *f.MyInfo))
}


func (f *FingerTable) FillTable() {
	/*
		Finds which baseHashGroups are populated (from Zookeeper)
		Choose successors to add to the finger table
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
	nrHashGroups := float64(len(baseHashGroupsInt))
	f.NrSuccessors = uint32((math.Log(nrHashGroups) / math.Log(2)) + 1)
	f.chooseSuccessors(baseHashGroupsInt)

	data, _ := json.Marshal(f.MyInfo)
	for baseHashGroup, _ := range f.Successors {
	_ = zkCli.RegisterEphemeralNode(zkCli.PrependFollowerPath(fmt.Sprintf("%d/",
														  baseHashGroup)), data)
	}
}


func (f *FingerTable) UpdateNodes() {
	/*
		Routine that constantly checks ChordUpdateChannel.ReqCh
		If there is a change in servers, the finger table is rebuild
	*/
	for {
		select {
		case <-util.ChordUpdateChannel.ReqCh:
			f.FillTable()
			go func() {
				util.ChordUpdateChannel.RespCh	<-true
			}()
		default:
		}
	}
}
