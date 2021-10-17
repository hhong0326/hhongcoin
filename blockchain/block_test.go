package blockchain

import (
	"reflect"
	"testing"

	"github.com/hhong0326/hhongcoin/utils"
)

func TestCreateBlock(t *testing.T) {

	dbStorage = fakeDB{}

	// add test mempool
	Mempool().Txs["test"] = &Tx{}
	// indepent with DB
	b := createBlock("x", 1, 1) // prevHash, height, diff

	if reflect.TypeOf(b) != reflect.TypeOf(&Block{}) {
		t.Error("createBlock() should return an interface of a block")
	}
}

// blockBytes := dbStorage.FindBlock(hash)
// 	if blockBytes == nil {
// 		return nil, ErrNotFound
// 	}
// 	block := &Block{}
// 	block.restore(blockBytes)
// 	return block, nil
func TestFindBlock(t *testing.T) {

	t.Run("Block not found", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				return nil
			},
		}

		_, err := FindBlock("x")
		if err == nil {
			t.Error("The block should not be found.")
		}
	})

	t.Run("Block is found", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				b := &Block{
					Height: 1,
				}
				return utils.ToBytes(b)
			},
		}

		block, _ := FindBlock("x")

		if block.Height != 1 {
			t.Error("Block should be found.")
		}
	})

}
