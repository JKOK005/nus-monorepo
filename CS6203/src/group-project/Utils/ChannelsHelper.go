package Utils

import (
	pb "group-project/Protobuf/Generate"
)

type getTermNoChannel struct {
	ReqCh 	chan bool
	RespCh 	chan uint32
} // Term number retrieving channel

type setTermNoChannel struct {
	ReqCh 	chan uint32
	RespCh 	chan bool
} 	// Term number setting channel

type setCycleNoChannel struct {
	ReqCh 	chan uint32
	RespCh 	chan bool
}

<<<<<<< HEAD
type putKeyChannel struct {
	ReqCh 	chan *pb.PutKeyMsg
	RespCh 	chan bool
}

type getKeyChannel struct {
	ReqCh 	chan string
	RespCh 	chan *pb.GetKeyResp
=======
type setRequestChannel struct {
	ReqCh	chan uint32
	RespCh	chan bool
}	// Data retrieval request channel (for clients)

type setPutChannel struct {
	ReqCh	chan uint32
	RespCh	chan bool
<<<<<<< HEAD
>>>>>>> finger table lookup complete
}
=======
}	// Data storing put channel (for clients)
>>>>>>> corrected fingertable

/*
Shared channels for all go routines to use for communication
*/
var (
<<<<<<< HEAD
	GetTermNoCh 	= &getTermNoChannel{ReqCh: make(chan bool), RespCh: make(chan uint32)}
	SetTermNoCh 	= &setTermNoChannel{ReqCh: make(chan uint32), RespCh: make(chan bool)}
	SetCycleNoCh 	= &setCycleNoChannel{ReqCh: make(chan uint32), RespCh: make(chan bool)}
	PutKeyChannel 	= &putKeyChannel{ReqCh: make(chan *pb.PutKeyMsg), RespCh: make(chan bool)}
	GetKeyChannel 	= &getKeyChannel{ReqCh: make(chan string), RespCh: make(chan *pb.GetKeyResp)}
)
=======
	GetTermNoCh 		= &getTermNoChannel{ReqCh: make(chan bool),
											RespCh: make(chan uint32)}
	SetTermNoCh 		= &setTermNoChannel{ReqCh: make(chan uint32),
											RespCh: make(chan bool)}
	SetCycleNoCh 		= &setCycleNoChannel{ReqCh: make(chan uint32),
											 RespCh: make(chan bool)}
	SetRequestChannel 	= &setRequestChannel{ReqCh: make(chan uint32),
											 RespCh: make(chan bool)}
	SetPutChannel 	 	= &setPutChannel{ReqCh: make(chan uint32),
										 RespCh: make(chan bool)}
)
>>>>>>> finger table lookup complete
