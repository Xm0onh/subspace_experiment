package operator

import (
	"fmt"
	"net/http"
	"reflect"

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
}
type operator struct {
	socket.Socket
	id          identity.NodeID
	MessageChan chan interface{}
	TxChan      chan interface{}
	txRange     int
	handles     map[string]reflect.Value
	mem         *mempool.Producer
	server      *http.Server
}

func NewOperator(id identity.NodeID) Operator {

	return &operator{
		id:          id,
		Socket:      socket.NewSocket(id, config.Configuration.Addrs),
		MessageChan: make(chan interface{}, 10240),
		TxChan:      make(chan interface{}, 10240),
		txRange:     500,
		handles:     make(map[string]reflect.Value),
		mem:         mempool.NewProducer(),
	}
}

func (o *operator) ID() identity.NodeID {
	return o.id
}

func (o *operator) Run() {
	log.Infof("node %v start running", o.id)
	// if len(o.handles) > 0 {
	// 	go o.handle()

	// }
	go o.recv()
	o.http()
}

func (o *operator) recv() {
	for {
		msg := o.Recv()
		fmt.Println("I am ", o.id, "and got the message", msg)
	}
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
		msg := <-o.MessageChan
		v := reflect.ValueOf(msg)
		name := v.Type().String()
		f, exists := o.handles[name]
		if !exists {
			log.Fatalf("no registered handle function for message type %v", name)
		}
		f.Call([]reflect.Value{v})
	}
}
