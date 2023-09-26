package blockchain

type Blockchain struct {
	longestTailBlock *Block
	view             int
}

func NewBlockchain() *Blockchain {
	bc := new(Blockchain)
	return bc
}

func (bc *Blockchain) AddBlock(block *Block) {
	// TODO
	bc.longestTailBlock = block
	bc.view = block.View

}
