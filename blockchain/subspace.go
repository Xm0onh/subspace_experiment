package blockchain

import (
	"fmt"

	"github.com/xm0onh/subspace_experiment/election"
	"github.com/xm0onh/subspace_experiment/identity"
	"github.com/xm0onh/subspace_experiment/log"
	"github.com/xm0onh/subspace_experiment/operator"
)

type Subspace struct {
	operator.Operator
	election.Election
	longestTailBlock chan *Block
	bc               *Blockchain
}

func NewSubpace(operator operator.Operator, elec election.Election, committedBlocks chan *Block) *Subspace {
	s := new(Subspace)
	s.Operator = operator
	s.Election = elec
	s.longestTailBlock = committedBlocks
	s.bc = NewBlockchain()
	return s
}

func (s *Subspace) ProcessBlock(proposer identity.NodeID, block *Block) error {
	if s.bc.view > block.View {
		return nil
	}
	// if s.FindLeaderFor(s.bc.view) != s.ID() {
	// 	return nil
	// }
	s.bc.AddBlock(block)
	s.bc.view++
	log.Debugf("New Block is Added. The Current View is %v", s.bc.view)
	// log.Debugf("Choosing new leader for view: %v", block.View+1)
	fmt.Println("Block Successfuly added")
	return nil
}

func (s *Subspace) GetView() int {
	return s.bc.view
}

func (s *Subspace) GetLeaderForFirstRound(view int) identity.NodeID {
	return s.FindLeaderFor(view)
}
