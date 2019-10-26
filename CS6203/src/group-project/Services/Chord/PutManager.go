package Chord


type PutManager struct {
	NodeAddr		string
	NodePort		uint32
	BaseHashGroup	uint32
	Protocol		*Protocol
}

func (p PutManager) Start() {
	p.Protocol, _ = NewProtocol(p.NodeAddr, p.NodePort, p.BaseHashGroup)

	// p.HandlePuts()
}
