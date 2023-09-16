package operator

import "github.com/xm0onh/subspace_experiment/mempool"

type Operator struct {
	txRange int
	mem     *mempool.Producer
}

func NewOperator() *Operator {
	o := new(Operator)

	// Must be set based on the formula
	o.txRange = 500
	o.mem = mempool.NewProducer()
	return o
}
