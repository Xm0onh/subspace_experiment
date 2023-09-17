package operator

import (
	"net/http"

	"github.com/xm0onh/subspace_experiment/identity"
	"github.com/xm0onh/subspace_experiment/mempool"
)

type Operator struct {
	id      identity.NodeID
	txRange int
	mem     *mempool.Producer
	server  *http.Server
}

func NewOperator(id identity.NodeID) *Operator {
	o := new(Operator)
	o.id = id
	// Must be set based on the formula
	o.txRange = 500
	o.mem = mempool.NewProducer()
	return o
}

func (o *Operator) Start() {
	o.http()
}
