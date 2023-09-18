package replica

import "github.com/xm0onh/subspace_experiment/blockchain"

type Inter interface {
	ProcessBlock(block *blockchain.Block) error
	GetView() int
}
