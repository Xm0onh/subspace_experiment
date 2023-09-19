package replica

import (
	"encoding/gob"
	"fmt"
	"sync/atomic"

	"github.com/xm0onh/subspace_experiment/blockchain"
	"github.com/xm0onh/subspace_experiment/config"
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
	r.Election = election.NewRotation(config.GetConfig().N())

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

func (r *Replica) processNewView(newView int) {

	log.Debugf("[%v] is processing new view: %v, leader is %v", r.ID(), newView, r.FindLeaderFor(newView))
	fmt.Println("the round is", newView, "The processor is", r.ID(), "the leader is", r.FindLeaderFor(newView), "the status is", r.IsLeader(r.ID(), newView))

	if !r.IsLeader(r.ID(), newView) {
		// fmt.Println("I am node,", r.ID(), "and I am not the leader for round", view)
		return
	}
	r.proposeBlock(newView)
}

func (r *Replica) proposeBlock(view int) {
	// r.totalBlockSize += len(block.Payload)
	block := blockchain.NewBlock(r.ID(), view, r.roundNo, r.roundNo-1, r.mem.GetTransactions())
	r.Broadcast(block)
	_ = r.Inter.ProcessBlock(r.ID(), block)
}

func (r *Replica) Start() {
	go r.Run()

	r.processNewView(0)

	<-r.start
	for r.isStarted.Load() {
		event := <-r.eventChan
		switch v := event.(type) {
		case blockchain.Block:
			r.Inter.ProcessBlock(r.ID(), &v)
			r.roundNo++
			r.processNewView(r.roundNo)
		}
	}

}
