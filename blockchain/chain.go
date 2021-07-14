package blockchain

import (
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
	b.CurrentDifficulty = block.Difficulty
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

func (b *blockchain) recalculateDifficulty() int {
	allBlocks := b.Blocks()
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

func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		// recal the difficulty
		return b.recalculateDifficulty()
	} else {
		return b.CurrentDifficulty
	}

}

// Singleton
func BlockChain() *blockchain {
	if b == nil { // only happend Once
		// Init
		once.Do(func() { // Only Once though there has many Goroutine starting
			b = &blockchain{ // Instance
				Height: 0,
			}
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
