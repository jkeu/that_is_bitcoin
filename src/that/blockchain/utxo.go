package blockchain

import "fmt"

type UTXO struct {
	Address  string
	TxID     string
	VOut     uint32
	Script   string
	Satoshis int64
}

func (my *UTXO) String() string {
	str := fmt.Sprintf("%s %d %d", my.TxID, my.VOut, my.Satoshis)
	return str
}
