package blockchain

import (
	"sync"
	"testing"

	"github.com/hhong0326/hhongcoin/utils"
)

type fakeDB struct {
	fakeLoadChain func() []byte
	fakeFindBlock func() []byte
}

// using fake if changes by IF
// need interface to get same type
func (f fakeDB) FindBlock(hash string) []byte {
	return f.fakeFindBlock()
}
func (f fakeDB) LoadChain() []byte {
	return f.fakeLoadChain()
}
func (f fakeDB) SaveBlock(hash string, data []byte) {}
func (f fakeDB) SaveChain(data []byte)              {}
func (f fakeDB) EmptyBlocks()                       {}

func TestBlockChain(t *testing.T) {

	t.Run("Should create blockchain", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeLoadChain: func() []byte {
				return nil
			},
		}

		bc := BlockChain()
		if bc.Height < 1 {
			t.Error("Blockchain() should create a blockchain")
		}
	})

	t.Run("Should restore blockchain", func(t *testing.T) {
		once = *new(sync.Once) // bcz shared variable in the same package
		dbStorage = fakeDB{
			fakeLoadChain: func() []byte {
				bc := &blockchain{
					NewestHash:        "x",
					Height:            2,
					CurrentDifficulty: 1,
				}

				return utils.ToBytes(bc)
			},
		}

		bc := BlockChain()
		if bc.Height != 2 {
			t.Error("Blockchain() should exist a blockchain")
		}
	})
}
