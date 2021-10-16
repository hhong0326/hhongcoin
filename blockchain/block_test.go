package blockchain

import (
	"reflect"
	"testing"
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
