package Utils

type TermNoChannel struct {
	ReqCh 	chan uint32
	RespCh 	chan bool
} 	// Term number setting channel