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

func (p *Producer) GetBundleCount() int {
	return p.mempool.actual_bundle
}

func (p *Producer) TxRangeAdjustment(expected_bundle, actual_bundle int) {
	p.mempool.TxRangeAdjustment(expected_bundle, actual_bundle)
}
