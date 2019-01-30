package utils

import (
	"crypto/ecdsa"
	"log"
)

// NetParams for net config
type NetParams struct {
	Legacy byte
	Segwit byte
	WIF    byte
	XPub   int32
	XPrv   int32
}

// ParamsMainNet initial mainnet
var ParamsMainNet = NetParams{
	Legacy: byte(0x00),
	Segwit: byte(0x05),
	WIF:    byte(0x80),
	XPub:   0x0488B21E,
	XPrv:   0x0488ADE4,
}

// ParamsTestNet initial testnet
var ParamsTestNet = NetParams{
	Legacy: byte(0x6f),
	Segwit: byte(0xc4),
	WIF:    byte(0xef),
	XPub:   0x043587CF,
	XPrv:   0x04358394,
}

// GenWallet generate address and private key
func GenWallet(testnet bool) {
	net := "mainnet"
	if testnet {
		net = "testnet"
	}
	log.Printf("Generate bitcoin %s address\n", net)

	myWallet := NewWallet()

	if !testnet {
		params := ParamsMainNet
		myMainNetLegacy := myWallet.GetAddressPubKey(&params)
		log.Println("Legacy:", string(myMainNetLegacy))
		myMainNetSegWit := myWallet.GetAddressScriptHash(&params)
		log.Println("Segwit:", string(myMainNetSegWit))
		myMainNetWif := myWallet.GetWif(&params)
		log.Println("Private key:", string(myMainNetWif))
	} else if testnet {
		params := ParamsTestNet
		myTestNetLegacy := myWallet.GetAddressPubKey(&params)
		log.Println("TestNet Legacy:", string(myTestNetLegacy))
		myTestNetSegwit := myWallet.GetAddressScriptHash(&params)
		log.Println("TestNet Segwit:", string(myTestNetSegwit))
		myTestNetWif := myWallet.GetWif(&params)
		log.Println("Private key:", string(myTestNetWif))
	}
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
	private, public := NewKeyPair(compress)
	wallet := Wallet{private, public, compress}

	return &wallet
}

// GetAddressPubKey gen bitcoin address
func (w Wallet) GetAddressPubKey(params *NetParams) []byte {
	pubKeyHash := Hash160(w.PublicKey)

	versionedPayload := append([]byte{params.Legacy}, pubKeyHash...)
	checksum := Checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}

// GetAddressScriptHash gen bitcoin address script hash
func (w Wallet) GetAddressScriptHash(params *NetParams) []byte {
	pubKeyHash := Hash160(w.PublicKey)
	op := []byte{byte(0x00), byte(0x14)}
	pubKeyHash = append(op, pubKeyHash...)
	pubKeyHash = Hash160(pubKeyHash)

	versionedPayload := append([]byte{params.Segwit}, pubKeyHash...)
	checksum := Checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}

// GetWif gen private key wif format
func (w Wallet) GetWif(params *NetParams) []byte {
	versionedPayload := append([]byte{params.WIF}, w.PrivateKey.D.Bytes()...)
	if w.Compress {
		versionedPayload = append(versionedPayload, byte(0x01))
	}
	checksum := Checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}
