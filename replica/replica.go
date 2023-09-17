package replica

import (
	"encoding/gob"
	"fmt"
	"sync/atomic"

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
}

func NewReplica(id identity.NodeID) *Replica {
	r := new(Replica)
	r.Operator = operator.NewOperator(id)

	//Election - TBD
	// r.Election = election.NewRotation(3, 0.1)

	r.mem = mempool.NewProducer()
	r.start = make(chan bool)
	r.eventChan = make(chan interface{})
	r.committedBlocks = make(chan *blockchain.Block, 100)
	gob.Register(blockchain.Block{})

	return r
}

func (r *Replica) startSignal() {
	fmt.Println("start signal")
	if !r.isStarted.Load() {
		fmt.Println("Started!")
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
	if !r.IsLeader(r.ID(), view) {
		return
	}
	r.proposeBlock(view)
}

func (r *Replica) proposeBlock(view int) {
	// r.totalBlockSize += len(block.Payload)
	block := blockchain.NewBlock(r.ID(), view, r.roundNo, r.roundNo-1, r.mem.GetTransactions())
	r.roundNo++
	r.Broadcast(block)
}

func (r *Replica) Start() {
	go r.Run()
	r.startSignal()
	<-r.start
	// go r.ListenLocalEvent()
	// go r.ListenCommittedBlocks()

	for r.isStarted.Load() {
		event := <-r.eventChan
		fmt.Println("hi", event)
		switch v := event.(type) {
		case int:
			r.processNewView(v)

		}
	}
}
