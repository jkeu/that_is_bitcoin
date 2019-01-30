package blockchain

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type ChainSoJson struct {
	Status string      `json:"status"`
	Data   ChainSoData `json:"data"`
}

type ChainSoData struct {
	Network string        `json:"network"`
	Address string        `json:"address"`
	Txs     []ChainSoUTXO `json:"txs"`
}

type ChainSoUTXO struct {
	TxID          string `json:"txid"`
	VOut          uint32 `json:"output_no"`
	ScriptAsm     string `json:"script_asm"`
	Script        string `json:"script_hex"`
	Amount        string `json:"value"`
	Confirmations int    `json:"confirmations"`
	TTime         int64  `json:"time"`
}

func GetUTXOChainSo(address string) ([]UTXO, error) {
	TestNetURL := "https://chain.so/api/v2/get_tx_unspent/BTCTEST/%s"
	MainNetURL := "https://chain.so/api/v2/get_tx_unspent/BTC/%s"

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

	// fmt.Println("## http get from", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var csJson ChainSoJson
	if err := json.Unmarshal(body, &csJson); err != nil {
		return nil, err
	}

	var dat = csJson.Data.Txs

	// fmt.Println(dat)

	sort.Slice(dat, func(i, j int) bool { return dat[i].Amount > dat[j].Amount })

	return chainSo2UTXO(address, dat), nil
}

func chainSo2UTXO(address string, chUTXO []ChainSoUTXO) []UTXO {
	dd := make([]UTXO, len(chUTXO))

	for i, ch := range chUTXO {
		dd[i].Address = address
		dd[i].TxID = ch.TxID
		dd[i].VOut = ch.VOut
		dd[i].Script = ch.Script

		var sat int64
		idx := strings.Index(ch.Amount, ".")
		newAm := ch.Amount[0:(idx-1)] + ch.Amount[(idx+1):]
		sat, err := strconv.ParseInt(newAm, 10, 64)
		if err != nil {
			return nil
		}

		dd[i].Satoshis = sat
	}

	return dd
}

// Broadcast to net
func Broadcast(hex string, testnet bool) error {
	TestNetURL := "https://chain.so/api/v2/send_tx/BTCTEST"
	MainNetURL := "https://chain.so/api/v2/send_tx/BTC"

	net := "mainnet"
	netURL := MainNetURL
	if testnet {
		net = "testnet"
		netURL = TestNetURL
	}

	fmt.Println("## http post ", netURL)

	resp, err := http.PostForm(netURL, url.Values{"tx_hex": {hex}})
	if err != nil {
		fmt.Println("http error:", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read from http error:", err)
		return err
	}

	fmt.Println(string(body))

	// TODO:
	fmt.Println("## broadcast to", net)
	return nil
}
