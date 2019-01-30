package utils

import (
	"fmt"
	"io/ioutil"
	"that/blockchain"
)

// Broadcast check balance from internet
func Broadcast(inFile string, testnet bool) {
	fmt.Println("### broadcast transaction ###")

	inBytes, err := ioutil.ReadFile(inFile)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println("## load transaction from", inFile)

	err = blockchain.Broadcast(string(inBytes), testnet)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("## done.\n")
}
