package blockchain

import (
	"reflect"
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

func TestBlocks(t *testing.T) {

	fakeBlocks := 0

	dbStorage = fakeDB{
		fakeFindBlock: func() []byte {
			var b *Block

			if fakeBlocks == 0 {
				b = &Block{
					Height:   2,
					PrevHash: "x",
				}
			}

			if fakeBlocks == 1 {
				b = &Block{
					Height: 1,
				}
			}

			fakeBlocks++
			return utils.ToBytes(b)
		},
	}

	bc := &blockchain{}
	blocks := Blocks(bc)

	if reflect.TypeOf(blocks) != reflect.TypeOf([]*Block{}) {
		t.Error("Blocks() should return a slice of blocks")
	}

}

func TestFindTx(t *testing.T) {

	t.Run("Tx not found", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {

				b := &Block{
					Height:       1,
					Transactions: []*Tx{},
				}
				return utils.ToBytes(b)
			},
		}
		tx := FindTx(&blockchain{NewestHash: "x"}, "test")

		if tx != nil {
			t.Error("FindTx() should be found none.")
		}
	})

	t.Run("Tx should be found", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {

				b := &Block{
					Height:       1,
					Transactions: []*Tx{{ID: "test"}},
				}
				return utils.ToBytes(b)
			},
		}
		tx := FindTx(&blockchain{NewestHash: "x"}, "test")

		if tx == nil {
			t.Error("FindTx() shound be found tx.")
		}
	})

}
