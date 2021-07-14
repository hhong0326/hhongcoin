package blockchain

import (
	"sync"

	"github.com/hhong0326/hhongcoin/db"
	"github.com/hhong0326/hhongcoin/utils"
)

type blockchain struct {
	// blocks []*Block // very long, then pointer!
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

// Singleton Pattern
// Want to be sharing only One Instance
var b *blockchain // == (b *blockchain) receiver
var once sync.Once

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) persist() {
	db.SaveCheckpoint(utils.ToBytes(b))
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func (b *blockchain) Blocks() []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash
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

// Singleton
func BlockChain() *blockchain {
	if b == nil { // only happend Once
		// Init
		once.Do(func() { // Only Once though there has many Goroutine starting
			b = &blockchain{"", 0} // Instance
			// search for checkpoint on the db
			if checkpoint := db.Checkpoint(); checkpoint == nil {
				b.AddBlock("Genesis")
			} else {
				// restore b from bytes
				b.restore(checkpoint)
			}
		}) // only one time call on go routine situation
	}
	return b // has been already Init
}
