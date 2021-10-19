package blockchain

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/hhong0326/hhongcoin/utils"
	"github.com/hhong0326/hhongcoin/wallet"
)

// transaction life-cycle
// mempool(Memory Pool) = 아직 확정되지 않은 거래 like array or slice before db

const (
	minerReward int = 50
)

type mempool struct {
	Txs map[string]*Tx
	m   sync.Mutex
}

// on memory
var m *mempool
var memOnce sync.Once

func Mempool() *mempool {
	memOnce.Do(func() {
		m = &mempool{
			Txs: make(map[string]*Tx),
		}
	})
	return m
}

type Tx struct {
	ID        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

type TxIn struct {
	TxID      string `json:"txId"`  // find previous txOuts
	Index     int    `json:"index"` // which txOut
	Signature string `json:"signature"`
}

type TxOut struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

type UTxOut struct {
	TxID   string `json:"txId"`
	Index  int    `json:"index"`
	Amount int    `json:"amount"`
}

func (tx *Tx) getId() {
	tx.ID = utils.Hash(tx)
}

func (tx *Tx) sign() {
	for _, txIn := range tx.TxIns {
		txIn.Signature = wallet.Sign(tx.ID, wallet.Wallet())
	}
}

// ***
func validate(tx *Tx) bool {
	valid := true
	for _, txIn := range tx.TxIns {
		prevTx := FindTx(BlockChain(), txIn.TxID)
		if prevTx == nil {
			valid = false
			break
		}
		address := prevTx.TxOuts[txIn.Index].Address
		valid = wallet.Verify(txIn.Signature, tx.ID, address)
		if !valid {
			break
		}
	}
	return valid
}

func isOnMempool(uTxOut *UTxOut) bool {
	exists := false
Outer:
	for _, tx := range Mempool().Txs {
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

var ErrorNoMoney = errors.New("not enough money")
var ErrorNotValid = errors.New("Tx Invalid")

// code challenge
func makeTx(from, to string, amount int) (*Tx, error) {
	if BalanceByAddress(from, BlockChain()) < amount {
		return nil, ErrorNoMoney
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
	tx.sign()
	valid := validate(tx)
	if !valid {
		return nil, ErrorNotValid
	}

	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) (*Tx, error) {
	m.m.Lock()
	defer m.m.Unlock()

	tx, err := makeTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		return nil, err
	}

	m.Txs[tx.ID] = tx

	return tx, nil
}

// chaining the tx to the block and making empty mempool
func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx(wallet.Wallet().Address)
	var txs []*Tx

	for _, tx := range m.Txs {
		txs = append(txs, tx)
	}

	txs = append(txs, coinbase)
	m.Txs = make(map[string]*Tx)

	return txs
}

func (m *mempool) AddPeerTx(tx *Tx) {
	m.m.Lock()
	defer m.m.Unlock()

	m.Txs[tx.ID] = tx
}

func GetMempool(m *mempool, rw http.ResponseWriter) {
	m.m.Lock()
	defer m.m.Unlock()

	utils.HandleErr(json.NewEncoder(rw).Encode(m.Txs))
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
