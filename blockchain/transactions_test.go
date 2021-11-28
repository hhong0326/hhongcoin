package blockchain

import (
	"net/http"
	"testing"
)

func TestSign(t *testing.T) {
	tx := &Tx{}

	t.Run("Test Sign", func(t *testing.T) {
		tx.sign()
	})
}

func TestValidate(t *testing.T) {
	// tx := &Tx{
	// 	ID:        "test",
	// 	Timestamp: 1,
	// 	TxIns: []*TxIn{
	// 		{
	// 			TxID:      "testTx",
	// 			Index:     0,
	// 			Signature: "",
	// 		},
	// 	},
	// 	TxOuts: []*TxOut{},
	// }

	// valid := validate(tx)
	// t.Log(valid)
	// if valid {
	// 	t.Error("validate() should return false")
	// }
}

func TestIsOnMempool(t *testing.T) {

}

// Mempool

func TestAddPeerTx(t *testing.T) {

	tx := &Tx{
		ID:        "test",
		Timestamp: 1,
		TxIns:     []*TxIn{},
		TxOuts:    []*TxOut{},
	}

	m.AddPeerTx(tx)

	if _, ok := m.Txs["test"]; !ok {
		t.Error("AddPeerTx() should return a tx")
	}
}

func TestGetMempool(t *testing.T) {

	var rw http.ResponseWriter
	GetMempool(Mempool(), rw)
}

// relative http
func TestAddTx(t *testing.T) {
	t.Run("Add Tx", func(t *testing.T) {
		// Mempool().Txs["test"] = &Tx{}
		// _, err := Mempool().AddTx("test", 100)

		// if err != nil {
		// 	t.Error("AddTx() should return a new tx")
		// }
	})
}

// cc
func TestMakeTx(t *testing.T) {

}
