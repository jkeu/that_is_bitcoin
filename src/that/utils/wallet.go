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
		myMainNetP2PKH := myWallet.GetAddressPubKeyHash(&params)
		log.Println("P2PKH:", string(myMainNetP2PKH))
		myMainNetP2SH := myWallet.GetAddressScriptHash(&params)
		log.Println("P2SH:", string(myMainNetP2SH))
		myMainNetP2WSH := myWallet.GetAddressWitnessScriptHash(&params)
		log.Println("P2WSH:", string(myMainNetP2WSH))
		myMainNetWif := myWallet.GetWif(&params)
		log.Println("Private key:", string(myMainNetWif))
	} else if testnet {
		params := ParamsTestNet
		myTestNetP2PKH := myWallet.GetAddressPubKeyHash(&params)
		log.Println("TestNet P2PKH:", string(myTestNetP2PKH))
		myTestNetP2SH := myWallet.GetAddressScriptHash(&params)
		log.Println("TestNet P2SH:", string(myTestNetP2SH))
		myTestNetP2WSH := myWallet.GetAddressWitnessScriptHash(&params)
		log.Println("TestNet P2WSH:", string(myTestNetP2WSH))
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

// GetAddressPubKeyHash gen bitcoin address : P2PKH
func (w Wallet) GetAddressPubKeyHash(params *NetParams) []byte {
	pubKeyHash := Hash160(w.PublicKey)

	versionedPayload := append([]byte{params.Legacy}, pubKeyHash...)
	checksum := Checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}

// GetAddressWitnessScriptHash gen bitcoin address script hash : P2WSH
// OP_0 PushData(hash160)
func (w Wallet) GetAddressWitnessScriptHash(params *NetParams) []byte {
	pubKeyHash := Hash160(w.PublicKey)
	op0 := []byte{byte(OpZero)}
	redeemScript := append(op0, PushData(pubKeyHash)...)
	addressHash := Hash160(redeemScript)

	versionedPayload := append([]byte{params.Segwit}, addressHash...)
	checksum := Checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}

// GetAddressScriptHash gen bitcoin address script hash : P2SH
// PushData(pubkey) OP_CHECKSIG
func (w Wallet) GetAddressScriptHash(params *NetParams) []byte {
	addressHash := Hash160(w.PayToScriptHashScript())

	versionedPayload := append([]byte{params.Segwit}, addressHash...)
	checksum := Checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}

// PayToScriptHashScript redeem script
func (w Wallet) PayToScriptHashScript() []byte {
	opCheckSig := []byte{byte(OpCheckSig)}
	redeemScript := append(PushData(w.PublicKey), opCheckSig...)

	return redeemScript
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
