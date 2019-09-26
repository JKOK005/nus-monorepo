package Utils

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

/*
 Shared channels for all go routines to use for communication
 */
var (
	GetTermNoCh 	= &getTermNoChannel{ReqCh: make(chan bool), RespCh: make(chan uint32)}
	SetTermNoCh 	= &setTermNoChannel{ReqCh: make(chan uint32), RespCh: make(chan bool)}
	SetCycleNoCh 	= &setCycleNoChannel{ReqCh: make(chan uint32), RespCh: make(chan bool)}
)