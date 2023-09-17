package blockchain

import (
	"github.com/xm0onh/subspace_experiment/identity"
	"github.com/xm0onh/subspace_experiment/mempool"
)

type Block struct {
	Proposer identity.NodeID
	// Sig
	view    int
	ID      int
	PrevID  int
	Payload []*mempool.Transaction
}

func NewBlock(proposer identity.NodeID, view, id, prevID int, payload []*mempool.Transaction) *Block {
	b := new(Block)
	b.Proposer = proposer
	b.view = view
	b.ID = id
	b.PrevID = prevID
	b.Payload = payload
	return b
}
