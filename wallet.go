package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"log"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"golang.org/x/crypto/ripemd160"
)

var bigRadix = big.NewInt(58)
var bigZero = big.NewInt(0)

const (
	versionLegacy      = byte(0x00)
	versionSegwit      = byte(0x05)
	addressChecksumLen = 4
	alphabet           = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	alphabetIdx0       = '1'
)

func main() {
	log.Println("Generate Bitcoin address")

	myWallet := NewWallet()

	myAddressLegacy := myWallet.GetAddressPublicKey()
	log.Println("Address Legacy:", string(myAddressLegacy))

	myAddressSegWit := myWallet.GetAddressScriptHash()
	log.Println("Address Segwit:", string(myAddressSegWit))

	myWif := myWallet.GetWif()
	log.Println("Private:", string(myWif))
}

// Wallet info
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
	Compress   bool
}

// NewWallet initial Wallet struct
func NewWallet() *Wallet {
	compress := true
	private, public := newKeyPair(compress)
	wallet := Wallet{private, public, compress}

	return &wallet
}

func newKeyPair(compress bool) (ecdsa.PrivateKey, []byte) {
	curve := btcec.S256()
	private, _ := ecdsa.GenerateKey(curve, rand.Reader)

	var pubKey []byte

	if !compress {
		format := byte(0x04)
		pubKey = append([]byte{format}, private.PublicKey.X.Bytes()...)
		pubKey = append(pubKey, private.PublicKey.Y.Bytes()...)
	} else {
		format := byte(0x02)
		if isOdd(private.PublicKey.Y) {
			format |= byte(0x01)
		}
		pubKey = append([]byte{format}, private.PublicKey.X.Bytes()...)
	}

	return *private, pubKey
}

// GetAddressPublicKey gen bitcoin address
func (w Wallet) GetAddressPublicKey() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{versionLegacy}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}

// GetAddressScriptHash gen bitcoin address script hash
func (w Wallet) GetAddressScriptHash() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)
	op := []byte{byte(0x00), byte(0x14)}
	pubKeyHash = append(op, pubKeyHash...)
	pubKeyHash = HashPubKey(pubKeyHash)

	versionedPayload := append([]byte{versionSegwit}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}

// GetWif gen private key wif format
func (w Wallet) GetWif() []byte {
	versionedPayload := append([]byte{byte(0x80)}, w.PrivateKey.D.Bytes()...)
	if w.Compress {
		versionedPayload = append(versionedPayload, byte(0x01))
	}
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}

// HashPubKey hash160(sha256(pubkey))
func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

func isOdd(a *big.Int) bool {
	return a.Bit(0) == 1
}

// Base58Encode encodes a byte slice to a modified base58 string.
func Base58Encode(b []byte) []byte {
	x := new(big.Int)
	x.SetBytes(b)

	answer := make([]byte, 0, len(b)*136/100)
	for x.Cmp(bigZero) > 0 {
		mod := new(big.Int)
		x.DivMod(x, bigRadix, mod)
		answer = append(answer, alphabet[mod.Int64()])
	}

	for _, i := range b {
		if i != 0 {
			break
		}
		answer = append(answer, alphabetIdx0)
	}

	alen := len(answer)
	for i := 0; i < alen/2; i++ {
		answer[i], answer[alen-1-i] = answer[alen-1-i], answer[i]
	}

	return answer
}
