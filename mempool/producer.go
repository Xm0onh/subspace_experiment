package mempool

type Producer struct {
	mempool *Mempool
}

func NewProducer() *Producer {
	return &Producer{
		mempool: NewMemPool(),
	}
}

func (p *Producer) GetTransactions() []Transaction {
	return p.mempool.SelectTransactionsForBundle()
}
