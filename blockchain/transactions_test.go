package blockchain

import (
	"testing"
)

func TestSign(t *testing.T) {

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

	// m.Txs["test"] = &Tx{}

	// var rw http.ResponseWriter
	// GetMempool(m, rw)
}
