package blockchain

type Blockchain struct {
	longestTailBlock    *Block
	highestCommitted    int
	committedBlockNo    int
	totalBlockIntervals int
	view                int
}

func NewBlockchain() *Blockchain {
	bc := new(Blockchain)
	return bc
}

func (bc *Blockchain) AddBlock(block *Block) {
	// TODO
	bc.longestTailBlock = block
	bc.highestCommitted = block.PrevID
	bc.committedBlockNo = block.ID
	bc.totalBlockIntervals++
	bc.view = block.View

}
