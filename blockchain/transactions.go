package blockchain

import (
	"errors"
	"time"

	"github.com/hhong0326/hhongcoin/utils"
)

// transaction life-cycle
// mempool(Memory Pool) = 아직 확정되지 않은 거래 like array or slice before db

const (
	minerReward int = 50
)

type mempool struct {
	Txs []*Tx
}

// on memory
var Mempool *mempool = &mempool{}

type Tx struct {
	ID        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

type TxIn struct {
	TxID  string `json:"txId"`  // find previous txOuts
	Index int    `json:"index"` // which txOut
	Owner string `json:"owner"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type UTxOut struct {
	TxID   string `json:"txId"`
	Index  int    `json:"index"`
	Amount int    `json:"amount"`
}

func (tx *Tx) getId() {
	tx.ID = utils.Hash(tx)
}

func isOnMempool(uTxOut *UTxOut) bool {
	exists := false
Outer:
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				exists = true
				break Outer
			}

		}
	}

	return exists
}

func makeCoinbaseTx(address string) *Tx {

	txIns := []*TxIn{
		{"", -1, "COINBASE"}, // Coinbase like printing currency
	}

	txOuts := []*TxOut{
		{address, minerReward},
	}

	tx := &Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}

	tx.getId()

	return tx
}

// func makeTx(from, to string, amount int) (*Tx, error) {
// if BlockChain().BalanceByAddress(from) < amount {
// 	return nil, errors.New("not enough money")
// }

// // create from txOuts
// var txIns []*TxIn
// var txOuts []*TxOut
// oldTxOuts := BlockChain().TxOutsByAddress(from)

// total := 0
// for _, txOut := range oldTxOuts {
// 	if total > amount {
// 		break
// 	}
// 	txIn := &TxIn{txOut.Owner, txOut.Amount}
// 	txIns = append(txIns, txIn)
// 	total += txOut.Amount
// }

// change := total - amount
// if change != 0 {
// 	changeTxOut := &TxOut{from, change}
// 	txOuts = append(txOuts, changeTxOut)
// }

// txOut := &TxOut{to, amount}
// txOuts = append(txOuts, txOut)

// tx := &Tx{
// 	Id:        "",
// 	Timestamp: int(time.Now().Unix()),
// 	TxIns:     txIns,
// 	TxOuts:    txOuts,
// }
// tx.getId()

// return tx, nil
// }

func makeTx(from, to string, amount int) (*Tx, error) {
	if BalanceByAddress(from, BlockChain()) < amount {
		return nil, errors.New("not enough money")
	}

	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0
	uTxOuts := UTxOutsByAddress(from, BlockChain())

	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break // 필요로 하는 tx input만
		}
		txIn := &TxIn{uTxOut.TxID, uTxOut.Index, from} // from is not secure
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}

	if change := total - amount; change != 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}

	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)

	tx := &Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()

	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("hong", to, amount)
	if err != nil {
		return err
	}

	m.Txs = append(m.Txs, tx)

	return nil
}

// chaining the tx to the block and making empty mempool
func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx("hong")
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil

	return txs
}

// Tx1
// 	TxIns[COINBASE]
// 	TxOuts[$5(you)] <---- Spent TxOut

// Tx2
// 	TxIns[Tx1.TxOuts[0]]
// 	TxOuts[$5(me)] <---- Spent TxOut for Tx3

// Tx3
// 	TxIns[Tx2.TxOuts[0]]
// 	TxOuts[$3(you), $2(me)] <---- Available uTxOut
