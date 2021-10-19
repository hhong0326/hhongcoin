package blockchain

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/hhong0326/hhongcoin/db"
	"github.com/hhong0326/hhongcoin/utils"
)

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
	blockInterval      int = 2
	allowedRange       int = 2
)

type blockchain struct {
	// blocks []*Block // very long, then pointer!
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
	m                 sync.Mutex
}

// interface for fake db testing
type storage interface {
	FindBlock(hash string) []byte
	LoadChain() []byte
	SaveBlock(hash string, data []byte)
	SaveChain(data []byte)
	EmptyBlocks()
}

// Singleton Pattern
// Want to be sharing only One Instance
var b *blockchain // == (b *blockchain) receiver
var once sync.Once

// to make variable
var dbStorage storage = db.DB{} // adapter

// Singleton
func BlockChain() *blockchain {
	// only happend Once
	// Init
	once.Do(func() { // Only Once though there has many Goroutine starting
		b = &blockchain{ // Instance
			Height: 0,
		}
		// search for checkpoint on the db
		if checkpoint := dbStorage.LoadChain(); checkpoint == nil {
			b.AddBlock()
		} else {
			// restore b from bytes
			b.restore(checkpoint)
		}
	}) // only one time call on go routine situation

	return b // has been already Init
}

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) AddBlock() *Block {
	block := createBlock(b.NewestHash, b.Height+1, getDifficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	persistBlockchain(b)

	return block
}

func persistBlockchain(b *blockchain) {
	dbStorage.SaveChain(utils.ToBytes(b))
}

// Any
func Blocks(b *blockchain) []*Block {
	b.m.Lock()
	defer b.m.Unlock()

	var blocks []*Block
	hashCursor := b.NewestHash // variable로 받아올 수 있다
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)

		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

func Txs(b *blockchain) []*Tx {
	var txs []*Tx
	for _, block := range Blocks(b) {
		txs = append(txs, block.Transactions...)
	}
	return txs
}

func FindTx(b *blockchain, targetID string) *Tx {
	for _, tx := range Txs(b) {
		if tx.ID == targetID {
			return tx
		}
	}
	return nil
}

// Any
func recalculateDifficulty(b *blockchain) int {
	allBlocks := Blocks(b)
	newestBlock := allBlocks[0]
	lastRecalculatedBlock := allBlocks[difficultyInterval-1]

	actualTime := (newestBlock.Timestamp / 60) - (lastRecalculatedBlock.Timestamp / 60)
	expectedTime := difficultyInterval * blockInterval

	if actualTime <= (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime >= (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1
	}

	return b.CurrentDifficulty
}

// Any
func getDifficulty(b *blockchain) int {

	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		// recal the difficulty
		return recalculateDifficulty(b)
	} else {
		return b.CurrentDifficulty
	}

}

// func (b *blockchain) txOuts() []*TxOut {

// 	var txOuts []*TxOut

// 	blocks := b.Blocks()

// 	for _, block := range blocks {
// 		for _, tx := range block.Transactions {
// 			txOuts = append(txOuts, tx.TxOuts...)
// 		}
// 	}

// 	return txOuts
// }

// api
// func (b *blockchain) TxOutsByAddress(address string) []*TxOut {

// 	var ownedTxOuts []*TxOut

// 	txOuts := b.txOuts()
// 	for _, txOut := range txOuts {
// 		if txOut.Owner == address {
// 			ownedTxOuts = append(ownedTxOuts, txOut)
// 		}
// 	}

// 	return ownedTxOuts
// }

// code challenge
// Unspent TxOuts
func UTxOutsByAddress(address string, b *blockchain) []*UTxOut {
	var uTxOuts []*UTxOut
	// spent txOuts
	creatorTxs := make(map[string]bool)

	for _, block := range Blocks(b) {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Signature == "COINBASE" {
					break
				}
				if FindTx(b, input.TxID).TxOuts[input.Index].Address == address {
					// I! this input can find tx what created txout
					// 사용자가 input 으로 사용하는 output을 찾아 그 output을 가진 tx의 id를 map에 저장
					// 이미 input으로 사용된 output을 소유한 txs 마킹
					creatorTxs[input.TxID] = true
				}
			}

			for i, output := range tx.TxOuts {
				if output.Address == address {
					if _, ok := creatorTxs[tx.ID]; !ok {
						//not found
						uTxOut := &UTxOut{tx.ID, i, output.Amount}
						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}

	return uTxOuts
}

func BalanceByAddress(address string, b *blockchain) int {

	ownedTxOuts := UTxOutsByAddress(address, b)
	var amount int

	for _, txOut := range ownedTxOuts {
		amount += txOut.Amount
	}

	return amount
}

func Status(b *blockchain, rw http.ResponseWriter) {
	b.m.Lock()
	defer b.m.Unlock()

	utils.HandleErr(json.NewEncoder(rw).Encode(b))
}

func (b *blockchain) Replace(newBlocks []*Block) {
	b.m.Lock()
	defer b.m.Unlock()

	b.CurrentDifficulty = newBlocks[0].Difficulty
	b.Height = len(newBlocks)
	b.NewestHash = newBlocks[0].Hash

	persistBlockchain(b)
	dbStorage.EmptyBlocks()

	for _, block := range newBlocks {
		persistBlock(block)
	}
}

func (b *blockchain) AddPeerBlock(newBlock *Block) {
	b.m.Lock()
	m.m.Lock()
	defer b.m.Unlock()
	defer m.m.Unlock()

	b.Height += 1
	b.CurrentDifficulty = newBlock.Difficulty
	b.NewestHash = newBlock.Hash

	persistBlockchain(b)
	persistBlock(newBlock)

	for _, tx := range newBlock.Transactions {
		_, ok := m.Txs[tx.ID]
		if ok {
			delete(m.Txs, tx.ID)
		}
	}
}

// What should be receiver / function
// flag : mutating struct except reading or not all time, Any blockchain | struct를 변화시키는 여부에 따라서
// If we are mutating the struct we will use receiver func.
// If we are using the struct as an input for data(reading) we will use a normal function.
