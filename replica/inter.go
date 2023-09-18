package replica

import (
	"github.com/xm0onh/subspace_experiment/blockchain"
	"github.com/xm0onh/subspace_experiment/identity"
)

type Inter interface {
	ProcessBlock(block *blockchain.Block) error
	GetView() int
	GetLeaderForFirstRound(view int) identity.NodeID
}
