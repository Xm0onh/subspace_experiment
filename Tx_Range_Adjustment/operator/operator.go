package operator

import (
	"fmt"
	"net"
	"reflect"

	"github.com/anthdm/hollywood/actor"

	"github.com/xm0onh/subspace_experiment/config"
	"github.com/xm0onh/subspace_experiment/identity"
	"github.com/xm0onh/subspace_experiment/log"
	"github.com/xm0onh/subspace_experiment/mempool"
	"github.com/xm0onh/subspace_experiment/socket"
)

type Operator interface {
	socket.Socket
	ID() identity.NodeID
	Run()
	Register(m interface{}, f interface{})
	RecvT() string
	SetT()
}
type operator struct {
	socket.Socket
	id          identity.NodeID
	MessageChan chan interface{}
	txRange     int
	handles     map[string]reflect.Value
	mem         *mempool.Producer
	test        string
}

func NewOperator(id identity.NodeID) Operator {

	return &operator{
		id:          id,
		Socket:      socket.NewSocket(id, config.Configuration.Addrs),
		MessageChan: make(chan interface{}, 10240),
		txRange:     500,
		handles:     make(map[string]reflect.Value),
		mem:         mempool.NewProducer(),
	}
}

func (o *operator) newServer(listenAddr string) actor.Producer {
	return func() actor.Receiver {
		return &server{
			listenAddr: listenAddr,
			sessions:   make(map[*actor.PID]net.Conn),
			operator:   o,
		}
	}
}

func (o *operator) newSession(conn net.Conn) actor.Producer {
	return func() actor.Receiver {
		return &session{
			conn:     conn,
			msg:      make(chan []byte, 1024),
			operator: o,
		}
	}
}

func (o *operator) newHandler() actor.Receiver {
	return &handler{}
}

func (o *operator) ID() identity.NodeID {
	return o.id
}

func (o *operator) Run() {
	log.Infof("node %v start running", o.id)
	if len(o.handles) > 0 {
		// go o.handle()
		// go o.recv()
	}

	o.http()
}

func (o *operator) Register(m interface{}, f interface{}) {
	t := reflect.TypeOf(m)
	fn := reflect.ValueOf(f)
	if fn.Kind() != reflect.Func {
		panic("handle function is not func")
	}

	if fn.Type().In(0) != t {
		panic("func type is not t")
	}

	if fn.Kind() != reflect.Func || fn.Type().NumIn() != 1 || fn.Type().In(0) != t {
		panic("register handle function error")
	}
	o.handles[t.String()] = fn
}

// handle receives messages from message channel and calls handle function using refection
func (o *operator) handle() {
	for {
		if o.test != "" {
			fmt.Println("there you go", o.test)
			o.test = ""
			// msg := blockchain.Block.FromString(o.test)
			// v := reflect.ValueOf(msg)
			// name := v.Type().String()
			// f, exists := o.handles[name]
			// if !exists {
			// 	log.Fatalf("no registered handle function for message type %v", name)
			// }

			// f.Call([]reflect.Value{v})
		}
	}
}

func (o *operator) RecvT() string {
	return o.test
}

func (o *operator) SetT() {
	o.test = ""
}
