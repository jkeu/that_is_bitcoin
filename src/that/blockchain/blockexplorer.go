package blockchain

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
)

// UTXO for unspent tx
type BlockExplorerUTXO struct {
	Address       string  `json:"address"`
	TxID          string  `json:"txid"`
	VOut          uint32  `json:"vout"`
	Script        string  `json:"scriptPubKey"`
	Amount        float64 `json:"amount"`
	Satoshis      int64   `json:"satoshis"`
	Height        int     `json:"height"`
	Confirmations int     `json:"confirmations"`
}

func GetUTXOBlockExplorer(address string) ([]UTXO, error) {
	TestNetURL := "https://testnet.blockexplorer.com/api/addr/%s/utxo"
	MainNetURL := "https://blockexplorer.com/api/addr/%s/utxo"

	var netURL string
	netFlag := address[0]

	if netFlag == '2' || netFlag == 'm' || netFlag == 'n' {
		netURL = TestNetURL
	} else if netFlag == '1' || netFlag == '3' {
		netURL = MainNetURL
	} else {
		err := errors.New("not support address type")
		return nil, err
	}

	url := fmt.Sprintf(netURL, address)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var dat []BlockExplorerUTXO
	if err := json.Unmarshal(body, &dat); err != nil {
		return nil, err
	}

	sort.Slice(dat, func(i, j int) bool { return dat[i].Satoshis > dat[j].Satoshis })

	return blockExplorer2UTXO(dat), nil
}

func blockExplorer2UTXO(beUTXO []BlockExplorerUTXO) []UTXO {
	dd := make([]UTXO, len(beUTXO))

	for i, be := range beUTXO {
		dd[i].TxID = be.TxID
		dd[i].VOut = be.VOut
		dd[i].Script = be.Script
		dd[i].Satoshis = be.Satoshis
	}

	return dd
}

func GetBalance(address string) (string, error) {
	TestNetURL := "https://testnet.blockexplorer.com/api/addr/%s/balance"
	MainNetURL := "https://blockexplorer.com/api/addr/%s/balance"

	var netURL string
	netFlag := address[0]

	if netFlag == '2' || netFlag == 'm' || netFlag == 'n' {
		netURL = TestNetURL
	} else if netFlag == '1' || netFlag == '3' {
		netURL = MainNetURL
	} else {
		err := errors.New("not support address type")
		return "", err
	}

	url := fmt.Sprintf(netURL, address)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
