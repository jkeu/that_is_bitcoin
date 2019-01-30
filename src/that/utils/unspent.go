package utils

import (
	"fmt"
	"that/blockchain"
)

// Unspent check utxo from internet
func Unspent(address string) {
	fmt.Println("### unspent list ###")

	utxo, err := blockchain.GetUTXOChainSo(address)
	if err != nil {
		panic(err)
	}

	for i, uu := range utxo {
		fmt.Printf(" *%d* %s %d %d\n", i, uu.TxID, uu.VOut, uu.Satoshis)
	}

	// fmt.Printf("%s %s %s\n", address, "value", balance)
}
