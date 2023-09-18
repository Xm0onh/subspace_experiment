package replica

import (
	"encoding/gob"
	"sync/atomic"
	"time"

	"github.com/xm0onh/subspace_experiment/blockchain"
	"github.com/xm0onh/subspace_experiment/election"
	"github.com/xm0onh/subspace_experiment/identity"
	"github.com/xm0onh/subspace_experiment/log"
	"github.com/xm0onh/subspace_experiment/mempool"
	"github.com/xm0onh/subspace_experiment/operator"
)

type Replica struct {
	operator.Operator
	election.Election
	mem              *mempool.Producer
	start            chan bool
	isStarted        atomic.Bool
	committedBlocks  chan *blockchain.Block
	eventChan        chan interface{}
	roundNo          int
	totalCommittedTx int
	totalBlockSize   int
	Inter
}

func NewReplica(id identity.NodeID) *Replica {
	r := new(Replica)
	r.Operator = operator.NewOperator(id)

	//Election - TBD
	r.Election = election.NewRotation(3, 0.3)

	r.mem = mempool.NewProducer()
	r.start = make(chan bool)
	r.eventChan = make(chan interface{})
	r.committedBlocks = make(chan *blockchain.Block, 1000)
	r.Register(blockchain.Block{}, r.HandleBlock)
	gob.Register(blockchain.Block{})

	r.Inter = blockchain.NewSubpace(r.Operator, r.Election, r.committedBlocks)
	return r
}

func (r *Replica) startSignal() {
	if !r.isStarted.Load() {
		r.isStarted.Store(true)
		log.Infof("Is Started = True")
		r.start <- true
	}
}

func (r *Replica) HandleBlock(block blockchain.Block) {
	r.startSignal()
	log.Debugf("[%v] received a block from %v, view is %v, id: %x, prevID: %x", r.ID(), block.Proposer, block.View, block.ID, block.PrevID)
	r.eventChan <- block
}
func (r *Replica) ListenLocalEvent() {

}

func (r *Replica) ListenCommittedBlocks() {

}

func (r *Replica) processNewView(view int) {
	log.Debugf("[%v] is processing new view: %v, leader is %v", r.ID(), view, r.FindLeaderFor(view))
	if r.GetLeader() != r.ID() {
		return
	}
	r.proposeBlock(view)
}

func (r *Replica) proposeBlock(view int) {
	// r.totalBlockSize += len(block.Payload)
	block := blockchain.NewBlock(r.ID(), view, r.roundNo, r.roundNo-1, r.mem.GetTransactions())
	r.Broadcast(block)
	_ = r.Inter.ProcessBlock(block)
}

func (r *Replica) Start() {
	go r.Run()
	// node_zero := identity.NewNodeID(1)
	// if r.ID() == node_zero {
	// 	r.proposeBlock(0)
	// }
	time.Sleep(1 * time.Second)
	l := r.Inter.GetLeaderForFirstRound(0)
	r.proposeBlock(l.Node())
	for {
		event := <-r.eventChan
		switch v := event.(type) {
		case blockchain.Block:
			r.processNewView(v.View + 1)
		}
	}

}
