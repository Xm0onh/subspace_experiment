package types

import (
	"github.com/xm0onh/subspace_experiment/identity"
	"github.com/xm0onh/subspace_experiment/mempool"
)

type Block struct {
	Proposer identity.NodeID
	// Sig
	View    int
	ID      int
	PrevID  int
	Payload []mempool.Transaction
}
