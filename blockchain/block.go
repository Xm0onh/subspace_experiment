package blockchain

import (
	"github.com/xm0onh/subspace_experiment/identity"
	"github.com/xm0onh/subspace_experiment/mempool"
)

type Block struct {
	Proposer identity.NodeID
	// Sig
	round   int
	ID      int
	PrevID  int
	Payload []*mempool.Transaction
}

func NewBlock(proposer identity.NodeID, round, id, prevID int, payload []*mempool.Transaction) *Block {
	b := new(Block)
	b.Proposer = proposer
	b.round = round
	b.ID = id
	b.PrevID = prevID
	b.Payload = payload
	return b
}
