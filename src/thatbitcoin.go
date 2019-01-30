package main

import (
	"flag"
	"fmt"
	"that/utils"
)

func usage() {
	fmt.Printf("Usage %s:\n", flag.Arg(0))
	fmt.Println("  thatbitcoin -balance -address (address)")
	fmt.Println("  thatbitcoin -unspent -address (address)")
	fmt.Println("  thatbitcoin -genkey [-testnet]")
	fmt.Println("  thatbitcoin -unsign -from (address) -to (address) -value (satoshis) -out (file)")
	fmt.Println("  thatbitcoin -sign -key (private key) -in (unsign file) -dump (signed file)")
	fmt.Println("  thatbitcoin -broadcast [-testnet] -file (signed file)")

	fmt.Printf("\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	balancePtr := flag.Bool("balance", false, "check address balance")
	addressPtr := flag.String("address", "", "bitcoin address")

	unspentPtr := flag.Bool("unspent", false, "get unspent list")

	genkeyPtr := flag.Bool("genkey", false, "generate a address and private key")
	testnetPtr := flag.Bool("testnet", false, "testnet")

	unsignPtr := flag.Bool("unsign", false, "create a unsign raw transaction")
	fromPtr := flag.String("from", "", "from address")
	toPtr := flag.String("to", "", "to address")
	valuePtr := flag.Int("value", 0, "transfer satoshis")
	outPtr := flag.String("out", "./unsign.txn", "output unsign file")

	signPtr := flag.Bool("sign", false, "sign a raw transaction")
	inPtr := flag.String("in", "./unsign.txn", "input unsign file")
	dumpPtr := flag.String("dump", "./signed.txn", "dump signed file")
	keyPtr := flag.String("key", "", "private key")

	broadcastPtr := flag.Bool("broadcast", false, "broadcast signed transaction")
	filePtr := flag.String("file", "./signed.txn", "signed file")

	flag.Parse()

	if *balancePtr {
		utils.Balance(*addressPtr)
	}

	if *unspentPtr {
		utils.Unspent(*addressPtr)
	}

	if *genkeyPtr {
		utils.GenWallet(*testnetPtr)
	}

	if *unsignPtr {
		if *fromPtr == "" || *toPtr == "" || *valuePtr == 0 {
			flag.Usage()
			return
		}

		utils.CreateUnsign(*fromPtr, *toPtr, int64(*valuePtr), *outPtr)
	}

	if *signPtr {
		if *keyPtr == "" {
			flag.Usage()
			return
		}
		utils.SignTx(*keyPtr, *inPtr, *dumpPtr)
	}

	if *broadcastPtr {
		utils.Broadcast(*filePtr, *testnetPtr)
	}
}
