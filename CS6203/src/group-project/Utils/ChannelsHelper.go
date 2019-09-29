package Utils

import (
	pb "group-project/Protobuf/Generate"
)

type getTermNoChannel struct {
	ReqCh 	chan bool
	RespCh 	chan uint32
} // Term number retriving channel

type setTermNoChannel struct {
	ReqCh 	chan uint32
	RespCh 	chan bool
} 	// Term number setting channel

type setCycleNoChannel struct {
	ReqCh 	chan uint32
	RespCh 	chan bool
}

type putKeyChannel struct {
	ReqCh 	chan *pb.PutKeyMsg
	RespCh 	chan *pb.PutKeyResp
}

type getKeyChannel struct {
	ReqCh 	chan *pb.GetKeyMsg
	RespCh 	chan *pb.GetKeyResp
}

/*
 Shared channels for all go routines to use for communication
 */
var (
	GetTermNoCh 	= &getTermNoChannel{ReqCh: make(chan bool), RespCh: make(chan uint32)}
	SetTermNoCh 	= &setTermNoChannel{ReqCh: make(chan uint32), RespCh: make(chan bool)}
	SetCycleNoCh 	= &setCycleNoChannel{ReqCh: make(chan uint32), RespCh: make(chan bool)}
	PutKeyChannel 	= &putKeyChannel{ReqCh: make(chan *pb.PutKeyMsg), RespCh: make(chan *pb.PutKeyResp)}
	GetKeyChannel 	= &getKeyChannel{ReqCh: make(chan *pb.GetKeyMsg), RespCh: make(chan *pb.GetKeyResp)}
)