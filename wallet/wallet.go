package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.com/hhong0326/hhongcoin/utils"
)

// 1) hash the msg
// " " -> hash(x) -> "hashed_message"
// 2) generate key pair
// Keypair (privateK, publicK) -> publicK will use in verify (Save priv to a file)
// 3) sign the hash
// ("hashed_message" + privateK) -> "signature"
// verify
// ("hashed_message" + "signature" + publicK) -> true / false

const (
	signature     string = "c5138ca9d6529cd3f71baf01c0026ad9dcb5a7a06dd64785dd11f83d8c883854e732d2aa5b627744e5f083240f628be379396a6d070f1ad8434dcf1f321ec7a5"
	privateKey    string = "30770201010420bff040c61409d2af2f28591faf960364d058bb013e92618840c1da92360f43d0a00a06082a8648ce3d030107a14403420004245b15e5eff0b28a422130deea9a0683461022447a5834742a0e738ad712550b1bdbb0850333bc3032cf43625a3beba5bcf9bd6a72b0cbb005d1164718d1020d"
	hashedMessage string = "c33084feaa65adbbbebd0c9bf292a26ffc6dea97b170d88e501ab4865591aafd"
)

const (
	fileName string = "hhongcoin.wallet"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

var w *wallet

func hasWalletFile() bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func createPrivKey() *ecdsa.PrivateKey {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)

	return privKey
}

func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)

	err = os.WriteFile(fileName, bytes, 0644) //read and write
	utils.HandleErr(err)
}

func restoreKey() (key *ecdsa.PrivateKey) { // return 될 variable을 사전에 선언
	keyAsBytes, err := os.ReadFile(fileName)
	utils.HandleErr(err)

	key, err = x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleErr(err)
	return // 필요없음
}

func encodeBigInts(a, b []byte) string {
	z := append(a, b...)
	return fmt.Sprintf("%x", z) // 16진수
}

func aFromK(key *ecdsa.PrivateKey) string {
	return encodeBigInts(key.X.Bytes(), key.Y.Bytes())
}

// hash msg + privateKey of wallet
func Sign(payload string, w *wallet) string {
	payloadAsBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)

	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsBytes)
	utils.HandleErr(err)

	return encodeBigInts(r.Bytes(), s.Bytes())
}

func restoreBigInts(payload string) (*big.Int, *big.Int, error) {
	bytes, err := hex.DecodeString(payload)
	if err != nil {
		return nil, nil, err
	}

	firstHalfBytes := bytes[:len(bytes)/2]
	secondHalfBytes := bytes[len(bytes)/2:]

	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(secondHalfBytes)
	return &bigA, &bigB, nil
}

// publickey + hash msg + signature
func Verify(signature, payload, address string) bool {

	r, s, err := restoreBigInts(signature)
	utils.HandleErr(err)
	x, y, err := restoreBigInts(address)
	utils.HandleErr(err)
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	payloadBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	ok := ecdsa.Verify(&publicKey, payloadBytes, r, s)
	return ok
}

// Singleton
func Wallet() *wallet {

	if w == nil {
		w = &wallet{}
		// has a wallet already?
		if hasWalletFile() {
			// yes -> restore from file
			w.privateKey = restoreKey()
		} else {
			// no -> create pre key, save to file
			key := createPrivKey()
			persistKey(key)
			w.privateKey = key
		}

		w.Address = aFromK(w.privateKey)
	}

	return w
}

func Start() {
	// privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	// utils.HandleErr(err)

	// keyAsBytes, err := x509.MarshalECPrivateKey(privateKey)
	// utils.HandleErr(err)

	// fmt.Printf("%x\n", keyAsBytes)

	// hashAsBytes, err := hex.DecodeString(hashedMessage)
	// utils.HandleErr(err)

	// r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashAsBytes) // divided signature r, s
	// utils.HandleErr(err)

	// signature := append(r.Bytes(), s.Bytes()...)
	// fmt.Printf("%x\n", signature)

	// ok := ecdsa.Verify(&privateKey.PublicKey, hashAsBytes, r, s)

	// fmt.Println(ok)

	// 16진수 체크
	privBytes, err := hex.DecodeString(privateKey)
	utils.HandleErr(err)

	private, err := x509.ParseECPrivateKey(privBytes)
	utils.HandleErr(err)

	fmt.Println(private)

	sigBytes, err := hex.DecodeString(signature)
	rBytes := sigBytes[:len(sigBytes)/2]
	sBytes := sigBytes[len(sigBytes)/2:]

	fmt.Printf("%d\n\n%d\n\n%d\n\n", sigBytes, rBytes, sBytes)

	var bigR, bigS = big.Int{}, big.Int{}

	bigR.SetBytes(rBytes)
	bigS.SetBytes(sBytes)

	hashBytes, err := hex.DecodeString(hashedMessage)
	utils.HandleErr(err)

	ok := ecdsa.Verify(&private.PublicKey, hashBytes, &bigR, &bigS)

	fmt.Println(ok)
}
