package utils

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"log"
	"math"
	"math/big"
	mrd "math/rand"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"golang.org/x/crypto/ripemd160"
)

var bigRadix = big.NewInt(58)
var bigZero = big.NewInt(0)

const (
	checksumLen  = 4
	alphabet     = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	alphabetIdx0 = '1'
)

var b58 = [256]byte{
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 0, 1, 2, 3, 4, 5, 6,
	7, 8, 255, 255, 255, 255, 255, 255,
	255, 9, 10, 11, 12, 13, 14, 15,
	16, 255, 17, 18, 19, 20, 21, 255,
	22, 23, 24, 25, 26, 27, 28, 29,
	30, 31, 32, 255, 255, 255, 255, 255,
	255, 33, 34, 35, 36, 37, 38, 39,
	40, 41, 42, 43, 255, 44, 45, 46,
	47, 48, 49, 50, 51, 52, 53, 54,
	55, 56, 57, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
}

var (
	// Used in RFC6979 implementation when testing the nonce for correctness
	one = big.NewInt(1)

	// oneInitializer is used to fill a byte slice with byte 0x01.  It is provided
	// here to avoid the need to create it multiple times.
	oneInitializer = []byte{0x01}
)

// NewKeyPair doc
func NewKeyPair(compress bool) (ecdsa.PrivateKey, []byte) {
	curve := btcec.S256()
	// curve := crypto.S256()
	private, _ := ecdsa.GenerateKey(curve, rand.Reader)

	var pubKey []byte

	if !compress {
		format := byte(0x04)
		pubKey = append([]byte{format}, private.PublicKey.X.Bytes()...)
		pubKey = append(pubKey, private.PublicKey.Y.Bytes()...)
	} else {
		format := byte(0x02)
		if IsOdd(private.PublicKey.Y) {
			format |= byte(0x01)
		}
		pubKey = append([]byte{format}, private.PublicKey.X.Bytes()...)
	}

	return *private, pubKey
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

// Base58Decode decode string to byte
func Base58Decode(b string) []byte {
	answer := big.NewInt(0)
	j := big.NewInt(1)

	scratch := new(big.Int)
	for i := len(b) - 1; i >= 0; i-- {
		tmp := b58[b[i]]
		if tmp == 255 {
			return []byte("")
		}
		scratch.SetInt64(int64(tmp))
		scratch.Mul(j, scratch)
		answer.Add(answer, scratch)
		j.Mul(j, bigRadix)
	}

	tmpval := answer.Bytes()

	var numZeros int
	for numZeros = 0; numZeros < len(b); numZeros++ {
		if b[numZeros] != alphabetIdx0 {
			break
		}
	}
	flen := numZeros + len(tmpval)
	val := make([]byte, flen)
	copy(val[numZeros:], tmpval)

	return val
}

// Hash160 hash160(sha256(pubkey))
func Hash160(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

// Checksum doc
func Checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:checksumLen]
}

// PushData len(data) + data
func PushData(data []byte) []byte {
	dlen := len(data)
	op := []byte{byte(dlen)}
	dataWithLen := append(op, data...)
	return dataWithLen
}

// IsOdd doc
func IsOdd(a *big.Int) bool {
	return a.Bit(0) == 1
}

// IsP2PKH for address
func IsP2PKH(address string) bool {
	if address[0] == '1' || address[0] == 'm' || address[0] == 'n' {
		return true
	}
	return false
}

// IsP2SH for address
func IsP2SH(address string) bool {
	if address[0] == '3' || address[0] == '2' {
		return true
	}
	return false
}

// GenerateNonce util
func GenerateNonce() [32]byte {
	var bytes [32]byte
	for i := 0; i < 32; i++ {
		//This is not "cryptographically random"
		bytes[i] = byte(randInt(0, math.MaxUint8))
	}
	return bytes
}

func randInt(min int, max int) uint8 {
	mrd.Seed(time.Now().UTC().UnixNano())
	return uint8(min + mrd.Intn(max-min))
}

// BytesEqual tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func BytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		// log.Printf("debug, len(a) %d, len(b) %d\n", len(a), len(b))
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
