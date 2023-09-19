package socket

import (
	"sync"
	"time"

	comm "github.com/xm0onh/subspace_experiment/communication"
	"github.com/xm0onh/subspace_experiment/identity"
	"github.com/xm0onh/subspace_experiment/log"
	"github.com/xm0onh/subspace_experiment/utils"
)

type Socket interface {
	Send(to identity.NodeID, m interface{})
	Broadcast(m interface{})
	Recv() interface{}
	Crash(t int)
	Close()
}
type socket struct {
	id        identity.NodeID
	addresses map[identity.NodeID]string
	crash     bool
	nodes     map[identity.NodeID]comm.IComm
	lock      sync.RWMutex
}

func NewSocket(id identity.NodeID, addrs map[identity.NodeID]string) Socket {

	socket := &socket{
		id:        id,
		addresses: addrs,
		crash:     false,
		nodes:     make(map[identity.NodeID]comm.IComm),
	}

	socket.nodes[id] = comm.NewComm(addrs[id])
	socket.nodes[id].Listen()

	return socket
}

func (s *socket) Send(to identity.NodeID, m interface{}) {
	s.lock.RLock()
	c, exists := s.nodes[to]
	s.lock.RUnlock()
	if !exists {
		s.lock.RLock()
		address, ok := s.addresses[to]
		s.lock.RUnlock()
		if !ok {
			log.Errorf("socket does not have address of node %s", to)
			return
		}
		c = comm.NewComm(address)
		err := utils.Retry(c.Dial, 100, time.Duration(50)*time.Millisecond)
		if err != nil {
			panic(err)
		}
		s.lock.Lock()
		s.nodes[to] = c
		s.lock.Unlock()
	}
	c.Send(s.id, m)
}

func (s *socket) Recv() interface{} {
	s.lock.RLock()
	c := s.nodes[s.id]
	s.lock.RUnlock()
	for {
		m := c.Recv()
		if !s.crash {
			return m
		}
	}
}

func (s *socket) Broadcast(m interface{}) {
	for id := range s.addresses {
		if id == s.id {
			s.Recv()
			continue
		}
		s.Send(id, m)
	}
}

func (s *socket) Close() {
	for _, c := range s.nodes {
		c.Close()
	}
}

func (s *socket) Crash(t int) {
	s.crash = true
	if t > 0 {
		timer := time.NewTimer(time.Duration(t) * time.Second)
		go func() {
			<-timer.C
			s.crash = false
		}()
	}
}
