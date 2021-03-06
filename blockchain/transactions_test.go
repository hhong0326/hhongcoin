package blockchain

import (
	"net/http/httptest"
	"testing"
)

type fakeMp struct {
	fakeMempool func()
}

// using fake if changes by IF
// need interface to get same type
func (f fakeMp) Mempool() {}

func TestSign(t *testing.T) {
	tx := &Tx{}

	t.Run("Test Sign", func(t *testing.T) {
		tx.sign()
	})
}

func TestValidate(t *testing.T) {

	// t.Run("Validate false: prevTx == nil", func(t *testing.T) {

	// 	once = *new(sync.Once) // bcz shared variable in the same package
	// 	dbStorage = fakeDB{
	// 		fakeFindBlock: func() []byte {

	// 			b := &Block{
	// 				Height: 1,
	// 				Hash:   "xxx",
	// 				Transactions: []*Tx{
	// 					{
	// 						ID: "test",
	// 						TxIns: []*TxIn{
	// 							{
	// 								TxID:      "test", // "test" of TestFindTx
	// 								Index:     0,
	// 								Signature: "xx",
	// 							},
	// 						},
	// 						TxOuts: []*TxOut{
	// 							{
	// 								Address: "x",
	// 								Amount:  0,
	// 							},
	// 						},
	// 					},
	// 				},
	// 			}

	// 			return utils.ToBytes(b)
	// 		},
	// 	}

	// 	valid := validate(Mempool().Txs["test"])
	// 	if valid {
	// 		t.Error("validate() should return false")
	// 	}

	// })

}

// Mempool

func TestIsOnMempool(t *testing.T) {

	t.Run("Not Exist", func(t *testing.T) {

		out := &UTxOut{
			TxID:   "test",
			Index:  0,
			Amount: 100,
		}

		if isOnMempool(out) {
			t.Error("isOnMempool should return false")
		}
	})

	t.Run("Is Exist", func(t *testing.T) {

		m := Mempool()
		m.Txs["test"] = &Tx{
			ID:        "test",
			Timestamp: 123,
			TxIns: []*TxIn{
				{
					TxID:      "test",
					Index:     0,
					Signature: "test",
				},
			},
			TxOuts: []*TxOut{
				{
					Address: "test",
					Amount:  100,
				},
			},
		}

		out := &UTxOut{
			TxID:   "test",
			Index:  0,
			Amount: 100,
		}

		if !isOnMempool(out) {
			t.Error("isOnMempool should return true")
		}
	})
}

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

	// var rw http.ResponseWriter

	rw := httptest.NewRecorder()
	GetMempool(Mempool(), rw)

	if rw.Result().StatusCode != 200 {
		t.Error("GetMempool should return 200")
	}
}

// relative http
func TestAddTx(t *testing.T) {
	t.Run("Add Tx", func(t *testing.T) {

		// tx := &Tx{
		// 	ID:        "test",
		// 	Timestamp: 1,
		// 	TxIns:     []*TxIn{},
		// 	TxOuts:    []*TxOut{},
		// }
		// _, err := Mempool().AddTx("test", 100)

		// if err != nil {
		// 	t.Error("AddTx() should return a new tx")
		// }
	})
}

// cc
func TestMakeTx(t *testing.T) {

}
