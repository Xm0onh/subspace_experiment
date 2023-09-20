package sendandreceive

import "github.com/xm0onh/subspace_experiment/identity"

type SandR struct {
	id  identity.NodeID
	msg chan []byte
}

func NewSandR(id identity.NodeID) *SandR {
	return &SandR{
		id:  id,
		msg: make(chan []byte),
	}
}

func (s *SandR) Send(msg []byte) {
	s.msg <- msg
}

func (s *SandR) Recv() []byte {
	return <-s.msg
}
