package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"that/blockchain"

	secp256k1 "github.com/toxeus/go-secp256k1"
)

const (
	// OpZero op_0
	OpZero = 0x0
	// OpDup opcode
	OpDup = 0x76
	// OpEqual opcode
	OpEqual = 0x87
	// OpEqualVerify opcode
	OpEqualVerify = 0x88
	// OpHash160 opcode
	OpHash160 = 0xa9
	// OpCheckSig opcode
	OpCheckSig = 0xac
)

const (
	txFee               = 1000
	defaultTxInOutAlloc = 15
)

// ThatTx top
type ThatTx struct {
	Version  int32
	TxIn     []*ThatTxIn
	TxOut    []*ThatTxOut
	LockTime uint32
}

// ThatTxIn in
type ThatTxIn struct {
	TxID      []byte
	Index     uint32
	SigScript []byte
	Witness   ThatWitness
	Sequence  uint32
}

// ThatWitness in
type ThatWitness [][]byte

// ThatTxOut out
type ThatTxOut struct {
	Value    int64
	PkScript []byte
}

// CreateUnsign create unsign raw transaction
func CreateUnsign(from, to string, value int64, outFile string) {
	fmt.Println("### create unsign transaction ###")

	unsignTx, err := createRawTx(from, to, value)
	if err != nil {
		fmt.Println(err)
		return
	}

	utx := unsignTx.Serialize()
	fmt.Println("## unsign hex:", hex.EncodeToString(utx))

	err = ioutil.WriteFile(outFile, []byte(hex.EncodeToString(utx)), 0644)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("## dump hex to", outFile)
}

// SignTx sign transaction
func SignTx(key, inFile, dumpFile string) {
	fmt.Println("### sign transaction ###")

	inBytes, err := ioutil.ReadFile(inFile)
	if err != nil {
		fmt.Print(err)
		return
	}
	unsignHex := string(inBytes) // load from unsign file

	unsignTx, err := decodeRawTx(unsignHex)
	if err != nil {
		log.Println("decodeRawTx error", err)
		return
	}
	// udeTx := unsignTx.Serialize()
	// fmt.Println("## decode data:", hex.EncodeToString(udeTx))

	signTx, err := signRawTx(unsignTx, key)
	if err != nil {
		log.Println(err)
		return
	}

	utx := signTx.Serialize()
	fmt.Println("## signed hex:", hex.EncodeToString(utx))

	err = ioutil.WriteFile(dumpFile, []byte(hex.EncodeToString(utx)), 0644)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("## dump hex to", dumpFile)
}

func decodeRawTx(utx string) (*ThatTx, error) {
	// log.Println("decodeRawTx entry", utx)
	//
	txBytes, err := hex.DecodeString(utx)
	if err != nil {
		newErr := fmt.Errorf("Decode error %s", err)
		return nil, newErr
	}
	idx := int(0)

	tx := newThatTx(0x0)

	versionBytes := make([]byte, 4)
	for i := 0; i < 4; i++ {
		versionBytes[i] = txBytes[i]
	}
	idx += 4

	r := bytes.NewReader(versionBytes)
	var version int32
	err = binary.Read(r, binary.LittleEndian, &version)
	if err != nil {
		newErr := fmt.Errorf("Decode version error %s", err)
		return nil, newErr
	}
	tx.Version = version

	marker := int8(txBytes[idx])
	idx++
	var inputs int8
	if marker == 0x0 {
		// TODO: witness version
		// read flag
		// read input
	} else {
		inputs = marker
	}

	var i int8
	for i = 0; i < inputs; i++ {
		// 32bytes hash
		hashSrcBytes := make([]byte, 32)
		for i := 0; i < 32; i++ {
			hashSrcBytes[i] = txBytes[idx+i]
		}

		hash := make([]byte, 32)
		for i := 0; i < len(hashSrcBytes); i++ {
			hash[i] = hashSrcBytes[len(hashSrcBytes)-i-1]
		}
		idx += 32
		// 4bytes index
		indexSrcBytes := make([]byte, 4)
		for i := 0; i < 4; i++ {
			indexSrcBytes[i] = txBytes[idx+i]
		}
		r = bytes.NewReader(indexSrcBytes)
		idx += 4
		var index uint32
		err = binary.Read(r, binary.LittleEndian, &index)
		if err != nil {
			newErr := fmt.Errorf("Decode in.index error %s", err)
			return nil, newErr
		}
		// sigScript
		scriptLen := int8(txBytes[idx])
		idx++
		script := make([]byte, scriptLen)
		for i := 0; i < int(scriptLen); i++ {
			script[i] = txBytes[idx+i]
		}
		idx += int(scriptLen)
		// sequence
		seqSrcBytes := make([]byte, 4)
		for i := 0; i < 4; i++ {
			seqSrcBytes[i] = txBytes[idx+i]
		}
		r = bytes.NewReader(seqSrcBytes)
		idx += 4
		var sequence uint32
		err = binary.Read(r, binary.LittleEndian, &sequence)
		if err != nil {
			newErr := fmt.Errorf("Decode in.sequence error %s", err)
			return nil, newErr
		}

		txIn, terr := newThatTxIn(hex.EncodeToString(hash), index, hex.EncodeToString(script), sequence)
		if terr != nil {
			newErr := fmt.Errorf("Decode new.in error %s", terr)
			return nil, newErr
		}

		tx.AddTxIn(txIn)
	}

	// outputs
	outputs := int8(txBytes[idx])
	idx++
	for i = 0; i < outputs; i++ {
		// value
		satSrcBytes := make([]byte, 8)
		for i := 0; i < 8; i++ {
			satSrcBytes[i] = txBytes[idx+i]
		}
		r = bytes.NewReader(satSrcBytes)
		idx += 8
		var satoshis int64
		err = binary.Read(r, binary.LittleEndian, &satoshis)
		if err != nil {
			newErr := fmt.Errorf("Decode out.value error %s", err)
			return nil, newErr
		}
		// script
		scriptLen := int8(txBytes[idx])
		idx++
		script := make([]byte, scriptLen)
		for i := 0; i < int(scriptLen); i++ {
			script[i] = txBytes[idx+i]
		}
		idx += int(scriptLen)

		var txOut ThatTxOut
		txOut.Value = satoshis
		txOut.PkScript = script
		tx.AddTxOut(&txOut)
	}

	return tx, nil
}

func signRawTx(tx *ThatTx, wif string) (*ThatTx, error) {
	secp256k1.Start()
	privateKeyBytes := Base58Decode(wif)
	privateKeyBytes = privateKeyBytes[1 : len(privateKeyBytes)-4]

	var privateKeyBytes32 [32]byte
	for i := 0; i < 32; i++ {
		privateKeyBytes32[i] = privateKeyBytes[i]
	}

	// Get the raw public key
	publicKeyBytes, success := secp256k1.Pubkey_create(privateKeyBytes32, true)
	if !success {
		log.Fatal("Failed to convert private key to public key")
	}

	// generate redeem script
	// mode-1 p2sh-segwit : OP_0 pushdata(hash160)
	pubKeyHash := Hash160(publicKeyBytes)
	Op0 := []byte{byte(OpZero)}
	redeemP2WSH := append(Op0, PushData(pubKeyHash)...)
	// mode-2 p2sh: pushdata(pubkey) OP_CHECKSIG
	opCheckSig := []byte{byte(OpCheckSig)}
	redeemP2SH := append(PushData(publicKeyBytes), opCheckSig...)
	redeemScript := redeemP2SH

	log.Println("# gen redeem script:", hex.EncodeToString(redeemScript))
	isP2SH := false

	for i, txIn := range tx.TxIn {
		sigScript := txIn.SigScript
		log.Println(i, "# txIn raw script:", hex.EncodeToString(sigScript))
		OpFirst := sigScript[0]
		if OpFirst == OpDup {
			log.Println("# I am P2PKH Transaction.")
			// P2PKH, script = OP_DUP OP_HASH160 PUSHDATA(address) OP_EQUALVERIFY OP_CHECKSIG
			break
		} else if OpFirst == OpHash160 {
			// P2SH, single script = PUSHDATA(pubKey) OP_CHECKSIG
			getHash := sigScript[2:22]
			wantHash := Hash160(redeemP2SH)
			// verify redeem script hash
			if !BytesEqual(getHash, wantHash) {
				log.Println("# I am P2WSH Transaction.")

				// common error
				newErr := fmt.Errorf("script signature not match, want %s but get %s",
					hex.EncodeToString(wantHash), hex.EncodeToString(getHash))

				// P2WSH not support yet
				// P2WSH redeem script = OP_0 PUSHDATA(address)
				wantHash = Hash160(redeemP2WSH)
				if BytesEqual(getHash, wantHash) {
					newErr = fmt.Errorf("Witness script signature not support yet")
				}

				return nil, newErr
			}

			log.Println("# I am P2SH Transaction.")
			isP2SH = true
			tx.TxIn[i].SigScript = redeemScript
			break
		} else {
			newErr := fmt.Errorf("sign script not support %s",
				hex.EncodeToString(sigScript))
			return nil, newErr
		}
	}

	rawTransaction := tx.Serialize()
	log.Println("debug: txIn:", hex.EncodeToString(tx.TxIn[0].SigScript))
	log.Println("debug: raw tx:", hex.EncodeToString(rawTransaction))

	// SIGHASH_ALL
	hashCodeType, err := hex.DecodeString("01000000")
	if err != nil {
		log.Fatal(err)
	}
	var rawTransactionBuffer bytes.Buffer
	rawTransactionBuffer.Write(rawTransaction)
	rawTransactionBuffer.Write(hashCodeType)
	rawTransaction = rawTransactionBuffer.Bytes()

	// Hash the raw transaction twice before the signing
	shaHash := sha256.New()
	shaHash.Write(rawTransaction)
	var hash []byte = shaHash.Sum(nil)

	shaHash2 := sha256.New()
	shaHash2.Write(hash)
	rawTransactionHashed := shaHash2.Sum(nil)
	var rawTransactionHashed32 [32]byte
	for i := 0; i < 32; i++ {
		rawTransactionHashed32[i] = rawTransactionHashed[i]
	}

	nonceByte := GenerateNonce()
	//Sign the raw transaction
	signedTransaction, success := secp256k1.Sign(rawTransactionHashed32, privateKeyBytes32, &nonceByte)
	if !success {
		log.Fatal("Failed to sign transaction")
	}

	//Verify that it worked.
	verified := secp256k1.Verify(rawTransactionHashed32, signedTransaction, publicKeyBytes)
	if !verified {
		log.Fatal("Failed to sign transaction")
	}

	secp256k1.Stop()

	hashCodeType, err = hex.DecodeString("01")
	if err != nil {
		log.Fatal(err)
	}

	// +1 for hashCodeType
	signedTransactionLength := byte(len(signedTransaction) + 1)

	var publicKeyBuffer bytes.Buffer
	publicKeyBuffer.Write(publicKeyBytes)
	pubKeyLength := byte(len(publicKeyBuffer.Bytes()))

	var buffer bytes.Buffer
	buffer.WriteByte(signedTransactionLength)
	buffer.Write(signedTransaction)
	buffer.WriteByte(hashCodeType[0])
	if !isP2SH {
		buffer.WriteByte(pubKeyLength)
		buffer.Write(publicKeyBuffer.Bytes())
	} else if isP2SH {
		// add redeemScript
		var redeemScriptBuffer bytes.Buffer
		redeemScriptBuffer.Write(redeemScript)
		redeemLength := byte(len(redeemScriptBuffer.Bytes()))
		buffer.WriteByte(redeemLength)
		buffer.Write(redeemScriptBuffer.Bytes())
	}

	scriptSig := buffer.Bytes()

	newTx := copyThatTx(tx)
	for i := 0; i < len(newTx.TxIn); i++ {
		newTx.TxIn[i].SigScript = scriptSig
	}

	// fmt.Println("### signScript data:", hex.EncodeToString(scriptSig))

	return newTx, nil
}

func createRawTx(from string, to string, satoshis int64) (*ThatTx, error) {
	/*
		if IsP2SH(from) {
			newErr := fmt.Errorf("%s SegWit-P2SH signature not support", from)
			return nil, newErr
		}
	*/
	utxos, err := blockchain.GetUTXOChainSo(from)
	if err != nil {
		newErr := fmt.Errorf("blockchain.GetUTXOChainSo error:%s", err)
		return nil, newErr
	}

	var balance int64
	for _, dd := range utxos {
		balance += dd.Satoshis
	}

	if balance < (satoshis + txFee) {
		newErr := fmt.Errorf("Error: Not enough satoshi, want %d but %d",
			satoshis+txFee, balance)
		return nil, newErr
	}

	tx := newThatTx(0x1)

	// inputs
	fmt.Println("## inputs:")
	balance = 0
	for i, uu := range utxos {
		// txIn, err := newThatTxIn(uu.TxID, uu.VOut, uu.Script, 0xfffffffd)
		txIn, terr := newThatTxIn(uu.TxID, uu.VOut, uu.Script, 0xffffffff)
		if terr != nil {
			return nil, terr
		}
		tx.AddTxIn(txIn)

		fmt.Printf(" *%d* %s\n", i, uu.String())

		balance += uu.Satoshis
		if balance >= (satoshis + txFee) {
			break
		}
	}

	// outputs
	fmt.Println("## outpus:")
	fmt.Printf(" *0* %35s %d\n", "[transaction fee]", txFee)
	txOut, err := newThatTxOut(satoshis, to)
	if err != nil {
		return nil, err
	}
	tx.AddTxOut(txOut)
	fmt.Printf(" *1* %35s %d\n", to, satoshis)

	// change
	change := balance - satoshis - txFee
	if change > 0 {
		txOut, err = newThatTxOut(change, from)
		if err != nil {
			return nil, err
		}
		tx.AddTxOut(txOut)
		fmt.Printf(" *2* %35s %d\n", from, change)
	}

	return tx, nil
}

func newThatTx(version int32) *ThatTx {
	return &ThatTx{
		Version: version,
		TxIn:    make([]*ThatTxIn, 0, defaultTxInOutAlloc),
		TxOut:   make([]*ThatTxOut, 0, defaultTxInOutAlloc),
	}
}

func copyThatTx(tx *ThatTx) *ThatTx {
	newTx := newThatTx(tx.Version)
	for _, tin := range tx.TxIn {
		newTx.AddTxIn(tin)
	}

	for _, tout := range tx.TxOut {
		newTx.AddTxOut(tout)
	}

	newTx.LockTime = tx.LockTime

	return newTx
}

// AddTxIn in
func (tx *ThatTx) AddTxIn(txIn *ThatTxIn) {
	tx.TxIn = append(tx.TxIn, txIn)
}

// AddTxOut out
func (tx *ThatTx) AddTxOut(txOut *ThatTxOut) {
	tx.TxOut = append(tx.TxOut, txOut)
}

// toJson for ThatTx
func (tx *ThatTx) toJSON() string {
	str := "{"

	str += fmt.Sprintf(`"version": %d, `, tx.Version)
	// inputs
	str += `"inputs": [`
	for _, txin := range tx.TxIn {
		str += fmt.Sprintf(`{"hash": "%s", "script": "%s", "index": %d}`,
			hex.EncodeToString(txin.TxID), hex.EncodeToString(txin.SigScript), txin.Index)
	}
	str += `], `
	// outputs
	str += `"outputs": [`
	for i, txout := range tx.TxOut {
		str += fmt.Sprintf(`{"script": "%s", "value": %d}`,
			hex.EncodeToString(txout.PkScript), txout.Value)
		if i < len(tx.TxOut)-1 {
			str += `, `
		}
	}
	str += `], `
	// locktime
	str += fmt.Sprintf(`"time": %d`, tx.LockTime)

	str += "}"

	return str
}

// Serialize for ThatTx
func (tx *ThatTx) Serialize() []byte {
	var buffer bytes.Buffer
	// version
	versionBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(versionBytes, uint32(tx.Version))
	buffer.Write(versionBytes)

	// TODO: marker & flag = "0001"
	// if witness add marker

	// input count
	inputs := len(tx.TxIn)
	buffer.WriteByte(byte(inputs))

	for _, in := range tx.TxIn {
		buffer.Write(in.TxID)

		indexBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(indexBytes, in.Index)
		buffer.Write(indexBytes)

		buffer.WriteByte(byte(len(in.SigScript)))
		buffer.Write(in.SigScript)

		sequenceBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(sequenceBytes, in.Sequence)
		buffer.Write(sequenceBytes)

		// log.Println("TxIn", i)
	}

	// output
	outputs := len(tx.TxOut)
	buffer.WriteByte(byte(outputs))

	for _, out := range tx.TxOut {
		satoshiBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(satoshiBytes, uint64(out.Value))
		buffer.Write(satoshiBytes)
		buffer.WriteByte(byte(len(out.PkScript)))
		buffer.Write(out.PkScript)

		// log.Println("TxOut", i)
	}

	timeBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(timeBytes, tx.LockTime)
	buffer.Write(timeBytes)

	return buffer.Bytes()
}

func newThatTxIn(txIDStr string, index uint32, script string, seq uint32) (*ThatTxIn, error) {
	hashBig, err := hex.DecodeString(txIDStr)
	if err != nil {
		return nil, err
	}
	hash := make([]byte, len(hashBig))
	for i := 0; i < len(hashBig); i++ {
		hash[i] = hashBig[len(hashBig)-i-1]
	}

	signatureScript, err := hex.DecodeString(script)
	if err != nil {
		return nil, err
	}

	var txIn ThatTxIn
	txIn.TxID = hash
	txIn.Index = index
	txIn.Sequence = seq
	txIn.SigScript = signatureScript

	return &txIn, nil
}

func newThatTxOut(satoshis int64, toAddress string) (*ThatTxOut, error) {
	// address := Base58.Decode(toAddress)
	address := Base58Decode(toAddress)
	address = address[1 : len(address)-4]
	// fmt.Println("address :", hex.EncodeToString(address))
	// fmt.Println("address1:", hex.EncodeToString(address1))

	var pkScript bytes.Buffer
	if IsP2SH(toAddress) {
		pkScript.WriteByte(byte(OpHash160))
		pkScript.WriteByte(byte(len(address))) //PUSH
		pkScript.Write(address)
		pkScript.WriteByte(byte(OpEqual))
	} else if IsP2PKH(toAddress) {
		pkScript.WriteByte(byte(OpDup))
		pkScript.WriteByte(byte(OpHash160))
		pkScript.WriteByte(byte(len(address))) //PUSH
		pkScript.Write(address)
		pkScript.WriteByte(byte(OpEqualVerify))
		pkScript.WriteByte(byte(OpCheckSig))
	} else {
		err := fmt.Errorf("Not support address type")
		return nil, err
	}

	var txOut ThatTxOut
	txOut.Value = satoshis
	txOut.PkScript = pkScript.Bytes()

	return &txOut, nil
}
