package socket

import (
	"sync"

	comm "github.com/xm0onh/subspace_experiment/communication"
	"github.com/xm0onh/subspace_experiment/identity"
	"github.com/xm0onh/subspace_experiment/log"
)

type socket struct {
	id        identity.NodeID
	addresses map[identity.NodeID]string
	nodes     map[identity.NodeID]comm.Comm

	lock sync.RWMutex
}

type Socket interface {
	Send(to identity.NodeID, msg interface{})

	// Broadcast send to all peers
	Broadcast(msg interface{})
	Recv() interface{}
	Close()
}

func NewSocket(id identity.NodeID, addrs map[identity.NodeID]string) Socket {

	socket := &socket{
		id:        id,
		addresses: addrs,
		nodes:     make(map[identity.NodeID]comm.Comm),
	}

	socket.nodes[id] = comm.NewComm(addrs[id])
	socket.nodes[id].Listen()

	return socket
}

func (s *socket) Send(to identity.NodeID, msg interface{}) {
	s.lock.RLock()
	t, exists := s.nodes[to]
	defer s.lock.RUnlock()
	if !exists {
		s.lock.RLock()
		address, ok := s.addresses[to]
		if !ok {
			log.Errorf("socket does not have address of node %s", to)
			return
		}
		t = comm.NewComm(address)
		s.lock.Lock()
		s.nodes[to] = t
		s.lock.Unlock()
	}
	t.Send(msg)
}

func (s *socket) Recv() interface{} {
	s.lock.RLock()
	t := s.nodes[s.id]
	s.lock.RUnlock()
	for {
		m := t.Recv()
		return m
	}
}

func (s *socket) Broadcast(m interface{}) {
	for id := range s.addresses {
		if id == s.id {
			continue
		}
		s.Send(id, m)
	}
}

func (s *socket) Close() {
	for _, t := range s.nodes {
		t.Close()
	}
}
