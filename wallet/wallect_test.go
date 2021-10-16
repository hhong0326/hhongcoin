package wallet

import (
	"crypto/x509"
	"encoding/hex"
	"io/fs"
	"reflect"
	"testing"
)

const (
	testKey     string = "30770201010420760e099316528ecbe23390e1feb36294b85f6eb823863f7f1a40395f5b451f93a00a06082a8648ce3d030107a1440342000468ff5e580e8957f71251d76abf4a7d4c1f2d21732e3f2a650477c738dddcc7f585fe0ff67076b1215a51b4940cf777ecc93612229938fe1d6493aba06d705b6b"
	testPayload string = "0fdb69f9c89d2033b3ccfe0af64c5d74492234b39df052d6c0f30d8b17372f5d"
	testSig     string = "091c543a83e44f26561dc57aefb12ee4b67a848c51cef2625ade3efe453871d14feede139728e820cdd5b4e7b3ade232c01db58880b3aed9e4344585e2e41045"
)

// interface 가지고 재구현
type fakeLayer struct {
	fakeHasWalletFile func() bool
}

// func (fakeLayer) hasWalletFile() bool {
// }

// ->
// (같은 func을 다른 방식으로)
func (f fakeLayer) hasWalletFile() bool {
	return f.fakeHasWalletFile()
}

func (fakeLayer) writeFile(name string, data []byte, perm fs.FileMode) error {
	return nil
}

func (fakeLayer) readFile(name string) ([]byte, error) {
	return x509.MarshalECPrivateKey(makeTestWallet().privateKey)
}
func TestWallet(t *testing.T) {

	t.Run("New Wallet is created", func(t *testing.T) {
		files = fakeLayer{ // (같은 함수를 다른 방식으로)
			fakeHasWalletFile: func() bool {
				t.Log("I have been called")
				return false
			},
		}

		tw := Wallet()
		if reflect.TypeOf(tw) != reflect.TypeOf(&wallet{}) {
			t.Error("New Wallet should return a new wallet instance")
		}

	})

	t.Run("Wallet is restored", func(t *testing.T) {
		files = fakeLayer{ // (같은 함수를 다른 방식으로)
			fakeHasWalletFile: func() bool {
				t.Log("I have been called")
				return true
			},
		}

		w = nil
		tw := Wallet()
		if reflect.TypeOf(tw) != reflect.TypeOf(&wallet{}) {
			t.Error("New Wallet should return a new wallet instance")
		}
	})
}

func makeTestWallet() *wallet {
	w := &wallet{}
	b, _ := hex.DecodeString(testKey)
	key, _ := x509.ParseECPrivateKey(b)
	w.privateKey = key
	w.Address = aFromK(key)

	return w
}

func TestSign(t *testing.T) {
	s := Sign(testPayload, makeTestWallet())
	_, err := hex.DecodeString(s)
	if err != nil {
		t.Errorf("Sign() should return a hex encoded string, got %s", s)
	}
	// t.Log(s)
}

func TestVerify(t *testing.T) {

	type test struct {
		input string
		ok    bool
	}

	tests := []test{
		{testPayload, true},
		{"0fdb69f9c89d2033b3ccfe0af64c5d74492234b39df052d6c0f30d8b17372f5f", false},
	}

	for _, tc := range tests {
		w := makeTestWallet()
		ok := Verify(testSig, tc.input, w.Address)

		if ok != tc.ok {
			t.Error("Verify() could not verify testSig and testPayload")
		}
	}
}

func TestRestoreBigInts(t *testing.T) {
	_, _, err := restoreBigInts("xx") // error moment
	if err == nil {
		t.Error("restoreBigInts() should return error when payload is not hex.")
	}
}
