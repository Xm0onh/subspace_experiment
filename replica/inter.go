package replica

import (
	"github.com/xm0onh/subspace_experiment/blockchain"
	"github.com/xm0onh/subspace_experiment/identity"
)

type Inter interface {
	ProcessBlock(proposer identity.NodeID, block *blockchain.Block) error
	GetView() int
	GetLeaderForFirstRound(view int) identity.NodeID
	AmIaLeader() identity.NodeID
}
