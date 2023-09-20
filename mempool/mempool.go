package mempool

import (
	"math"
	"math/rand"
)

type Mempool struct {
	txns             []Transaction
	uMax             uint64
	expected_bundle  int
	actual_bundle    int
	current_tx_range uint64
}

type Transaction struct {
	From   string
	To     string
	Amount int
	Hash   string
	Seq    int
}

const (
	lettersAndNumbers         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	INITIAL_TX_RANGE_DIVISOR  = 3 // U256::MAX/3
	TX_RANGE_ADJUSTMENT_RATIO = 4 // maximum or minimum ratio for adjustment
	MAX                       = 2 ^ 20
)

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
		txns:             transactionGenerator(int(MAX / 3)),
		uMax:             MAX,
		current_tx_range: MAX,
	}
}

func (m *Mempool) BidirectionalDistance(hash1, hash2 string) int {
	// TODO:
	return 0
}

func (m *Mempool) SelectTransactionsForBundle(vrfSignature string) []Transaction {
	// slot_vrf_hash := hash(vrfSignature)  // TODO: implement a hash function
	selectedTxs := []Transaction{}

	// for _, tx := range m.txns {
	// 	distance := m.BidirectionalDistance("0", tx.Hash)
	// 	if distance <= m.current_tx_range/2 {
	// 		selectedTxs = append(selectedTxs, tx)
	// 	}
	// }

	return selectedTxs
}

func (m *Mempool) GetTx() []Transaction {
	return m.txns[0:2]
}

func (m *Mempool) TxRangeAdjustment() {
	newRange := math.Max(
		math.Min(float64(m.expected_bundle)/float64(m.actual_bundle), 4), 0.25) * float64(m.current_tx_range)
	m.current_tx_range = uint64(newRange)
}
