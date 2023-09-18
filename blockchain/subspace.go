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
	committedBlocks chan *Block
	bc              *Blockchain
}

func NewSubpace(operator operator.Operator, elec election.Election, committedBlocks chan *Block) *Subspace {
	s := new(Subspace)
	s.Operator = operator
	s.Election = elec
	s.committedBlocks = committedBlocks
	s.bc = NewBlockchain()
	return s
}

func (s *Subspace) ProcessBlock(proposer identity.NodeID, block *Block) error {
	if s.bc.view > block.View {
		return nil
	}
	lead := s.FindLeaderFor(s.bc.view)
	fmt.Println("s.bc.view", s.bc.view, "leader is", lead)
	if (s.Election.IsLeader(proposer, s.bc.view)) && (s.bc.view != 0) {
		return nil
	}
	s.bc.AddBlock(block)
	log.Debugf("New Block is Added. The Current View is %v", s.bc.view)
	// log.Debugf("Choosing new leader for view: %v", block.View+1)
	return nil
}

func (s *Subspace) GetView() int {
	return s.bc.view
}

func (s *Subspace) GetLeaderForFirstRound(view int) identity.NodeID {
	return s.FindLeaderFor(view)
}
