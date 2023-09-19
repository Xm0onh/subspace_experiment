package socket

import (
	"net/url"
	"sync"

	comm "github.com/xm0onh/subspace_experiment/communication"
	"github.com/xm0onh/subspace_experiment/config"
	"github.com/xm0onh/subspace_experiment/identity"
	"github.com/xm0onh/subspace_experiment/log"

	"fmt"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/examples/remote/msg"
	"github.com/anthdm/hollywood/remote"
)

type Socket interface {
	// Send(to identity.NodeID, m interface{})
	Send(to identity.NodeID, msg string)
	// Broadcast(m interface{})
	Broadcast(msg string)
	Recv() interface{}
	Close()
	connect(port string)
}
type socket struct {
	id        identity.NodeID
	addresses map[identity.NodeID]string
	crash     bool
	nodes     map[identity.NodeID]comm.IComm
	lock      sync.RWMutex
	cPort     string
	// sPort     string
	e *actor.Engine
	r *remote.Remote
}

func NewSocket(id identity.NodeID, addrs map[identity.NodeID]string) Socket {

	cIp, cErr := url.Parse(config.Configuration.Addrs[id])
	if cErr != nil {
		log.Fatal("http url parse error: ", cErr)
	}
	cPort := ":" + cIp.Port()
	socket := &socket{
		id:        id,
		addresses: addrs,
		crash:     false,
		nodes:     make(map[identity.NodeID]comm.IComm),
		cPort:     cPort,
	}
	socket.connect(cPort)
	// connect(cPort, sPort)
	// socket.nodes[id] = comm.NewComm(addrs[id])
	// socket.nodes[id].Listen()

	return socket
}

func (s *socket) connect(port string) {
	fmt.Println("127.0.0.1" + port)
	// fmt.Println("127.0.0.1" + sPort)
	e := actor.NewEngine()
	r := remote.New(e, remote.Config{ListenAddr: "127.0.0.1" + (port)})
	e.WithRemote(r)
	s.e = e
	s.r = r
	// pid := actor.NewPID("127.0.0.1"+sPort, "server")
	// return pid

	// for {
	// 	e.Send(pid, &msg.Message{Data: "hello!" + sPort})
	// 	time.Sleep(2 * time.Second)
	// }
}

func (s *socket) Send(to identity.NodeID, m string) {
	sIp, sErr := url.Parse(config.Configuration.HTTPAddrs[to])
	if sErr != nil {
		log.Fatal("http url parse error: ", sErr)
	}
	sPort := ":" + sIp.Port()
	pid := actor.NewPID("127.0.0.1"+sPort, "server")
	s.e.Send(pid, &msg.Message{Data: m})

}

// func (s *socket) Send(to identity.NodeID, m interface{}) {
// 	s.lock.RLock()
// 	c, exists := s.nodes[to]
// 	s.lock.RUnlock()
// 	if (!exists) || to == s.id {
// 		s.lock.RLock()
// 		address, ok := s.addresses[to]
// 		s.lock.RUnlock()
// 		if !ok {
// 			log.Errorf("socket does not have address of node %s", to)
// 			return
// 		}
// 		c = comm.NewComm(address)
// 		err := utils.Retry(c.Dial, 100, time.Duration(50)*time.Millisecond)
// 		if err != nil {
// 			panic(err)
// 		}
// 		s.lock.Lock()
// 		s.nodes[to] = c
// 		s.lock.Unlock()
// 	}

// 	c.Send(s.id, m)

// }

func (s *socket) Recv() interface{} {
	s.lock.RLock()
	c := s.nodes[s.id]
	s.lock.RUnlock()
	for {
		m := c.Recv()
		return m

	}
}

func (s *socket) Broadcast(msg string) {
	for id := range s.addresses {
		s.Send(id, msg)
	}
}

func (s *socket) Close() {
	for _, c := range s.nodes {
		c.Close()
	}
}
