package utils

import (
	"fmt"
	"that/blockchain"
)

// Balance check balance from internet
func Balance(address string) {
	fmt.Println("### get balance ###")

	balance, err := blockchain.GetBalance(address)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s %s %s\n", address, "value", balance)
}
