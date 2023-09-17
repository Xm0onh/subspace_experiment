package mempool

import (
	"math/rand"
)

type Mempool struct {
	txns []Transaction
}

type Transaction struct {
	From   string
	To     string
	Amount int
	Hash   string
	Seq    int
}

const lettersAndNumbers = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = lettersAndNumbers[rand.Intn(len(lettersAndNumbers))]
	}
	return string(b)
}

func generateTransaction(seq int) Transaction {
	randStringBytes(35)
	randStringBytes(35)
	rand.Intn(100000000000000)
	return Transaction{
		//	From:   "1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2",
		//To:     "3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy",
		From:   randStringBytes(35),
		To:     randStringBytes(35),
		Amount: 1,
		Hash:   randStringBytes(35),
		Seq:    seq,
	}
}

func transactionGenerator(numberOfTransactions int) []Transaction {
	transactionList := make([]Transaction, 0)
	for i := 0; i < numberOfTransactions; i++ {
		transactionList = append(transactionList, generateTransaction(i))
	}
	return transactionList
}

func NewMemPool() *Mempool {
	return &Mempool{
		txns: transactionGenerator(500),
	}
}

func (m *Mempool) GetTx() []Transaction {
	return m.txns[0:2]
}
