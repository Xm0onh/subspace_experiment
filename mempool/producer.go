package mempool

type Producer struct {
	mempool *Mempool
}

func NewProducer() *Producer {
	return &Producer{
		mempool: NewMemPool(),
	}
}
