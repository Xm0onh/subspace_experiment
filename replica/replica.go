package replica

import (
	"encoding/gob"
	"fmt"
	"time"

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
	mem             *mempool.Producer
	committedBlocks chan *blockchain.Block
	eventChan       chan interface{}
	roundNo         int
	Inter
}

func NewReplica(id identity.NodeID) *Replica {
	r := new(Replica)
	r.Operator = operator.NewOperator(id)
	r.Election = election.NewRotation(config.GetConfig().N())
	r.mem = mempool.NewProducer()
	r.roundNo = 0
	r.eventChan = make(chan interface{})
	r.committedBlocks = make(chan *blockchain.Block, 1000)
	r.Register(blockchain.Block{}, r.HandleBlock)
	gob.Register(blockchain.Block{})

	r.Inter = blockchain.NewSubpace(r.Operator, r.Election, r.committedBlocks)
	return r
}

func (r *Replica) HandleBlock(block blockchain.Block) {
	if r.IsLeader(r.ID(), r.roundNo) {
		fmt.Println("Diiibs")
		_ = r.Inter.ProcessBlock(r.ID(), &block)
	}
	// if !r.IsLeader(r.ID(), r.roundNo) {
	// 	r.roundNo++
	// }
	log.Debugf("[%v] received a block from %v, view is %v, id: %x, prevID: %x", r.ID(), block.Proposer, block.View, block.ID, block.PrevID)
	// r.roundNo++
	r.eventChan <- block

}

func (r *Replica) proposeBlock(view int) {

	// log.Debugf("[%v] is processing new view: %v, leader is %v", r.ID(), view, r.FindLeaderFor(view))
	block := blockchain.NewBlock(r.ID(), view, r.roundNo, r.roundNo-1, r.mem.GetTransactions())
	r.Broadcast(block)
	// time.Sleep(300 * time.Millisecond)
}

func (r *Replica) Start() {
	go r.Run()

	if r.IsLeader(r.ID(), r.roundNo) && r.roundNo == 0 {
		r.proposeBlock(r.roundNo)
		fmt.Println("Hello")
	}
	for {

		event := <-r.eventChan
		if block, ok := event.(blockchain.Block); ok {
			_ = r.Inter.ProcessBlock(r.ID(), &block)
			// r.Inter.ProcessBlock(r.ID(), &block)
			fmt.Println("I am node", r.ID(), "and I received a block in round", r.roundNo-1)
			fmt.Println("Leader for", r.roundNo, "is:-->", r.FindLeaderFor(r.roundNo))
			fmt.Println("----My view is", r.roundNo-1)
			fmt.Println()
			r.roundNo++

			startTime := time.Now()
			for {
				// Check if 200 milliseconds have passed
				if time.Since(startTime) >= 500*time.Millisecond {
					r.roundNo++
					r.proposeBlock(r.roundNo)
					break
				}
				r.proposeBlock(r.roundNo)
			}

		}
	}

}
