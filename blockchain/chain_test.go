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

	blocks := []*Block{
		{Height: 2, PrevHash: "x"},
		{Height: 1, PrevHash: ""}, // Genesis block
	}

	fakeBlock := 0

	dbStorage = fakeDB{
		fakeFindBlock: func() []byte {

			defer func() {
				fakeBlock++
			}()
			return utils.ToBytes(blocks[fakeBlock])
		},
	}

	bc := &blockchain{}
	blocksResult := Blocks(bc)

	if reflect.TypeOf(blocksResult) != reflect.TypeOf([]*Block{}) {
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

func TestReplace(t *testing.T) {
	bc := &blockchain{
		Height:            1,
		CurrentDifficulty: 1,
		NewestHash:        "x",
	}

	blocks := []*Block{
		{Difficulty: 2, Hash: "xx"},
		{Difficulty: 2, Hash: "xx"},
	}

	bc.Replace(blocks)

	if bc.Height != 2 || bc.CurrentDifficulty != 2 || bc.NewestHash != "xx" {
		t.Error("Replace() should mutate the blockchain")
	}
}

func TestAddPeerBlock(t *testing.T) {

	bc := &blockchain{
		Height:            1,
		CurrentDifficulty: 1,
		NewestHash:        "xx",
	}

	m.Txs["test"] = &Tx{}

	newb := &Block{
		Difficulty: 2,
		Hash:       "test",
		Transactions: []*Tx{
			{ID: "test"},
		},
	}

	bc.AddPeerBlock(newb)

	if bc.Height != 2 || bc.CurrentDifficulty != 2 || bc.NewestHash != "test" {
		t.Error("AddPeerBlock() should mutate the blockchain.")
	}
}

func TestGetDifficulty(t *testing.T) {

	blocks := []*Block{
		{PrevHash: "x"},
		{PrevHash: "x"},
		{PrevHash: "x"},
		{PrevHash: "x"},
		{PrevHash: ""},
	}

	fakeBlock := 0
	dbStorage = fakeDB{
		fakeFindBlock: func() []byte {
			defer func() {
				fakeBlock++
			}()
			return utils.ToBytes(blocks[fakeBlock])
		},
	}

	type test struct {
		height int
		want   int
	}

	tests := []test{
		{height: 0, want: defaultDifficulty},
		{height: 2, want: defaultDifficulty},
		{height: 5, want: defaultDifficulty + 1},
	}

	for _, tc := range tests {
		bc := &blockchain{Height: tc.height, CurrentDifficulty: defaultDifficulty}
		got := getDifficulty(bc)

		if got != tc.want {
			t.Errorf("getDifficulty() should return %d, got %d", tc.want, got)
		}
	}
}

func TestUTxOutsByAddress(t *testing.T) {

	bc := &blockchain{}
	address := "x"

	utxOuts := UTxOutsByAddress(address, bc)

	if len(utxOuts) == 0 {
		t.Error("UTxOutsByAddress() should return result")
	}
}

func TestBalanceByAddress(t *testing.T) {
}
