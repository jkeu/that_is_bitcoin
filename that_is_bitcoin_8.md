# Python -  比特幣未簽名交易

[上篇](https://www.facebook.com/notes/%E5%82%91%E5%85%8B%E7%A7%91%E6%8A%80%E9%A4%A8/python-%E7%94%9F%E6%88%90%E6%AF%94%E7%89%B9%E5%B9%A3%E5%9C%B0%E5%9D%80/2315035778549174/)我們生成了地址，有了地址意味著你可以收幣，接下來我們來看看怎麼轉帳，有入有出才算真正學會使用比特幣。

交易簽名跟地址類型有關係，因為轉出就需要解碼上一筆的轉入交易。

簽名過程我們分為兩步，先在 online 電腦上生成未簽名交易，然後把數據複製到 offline 電腦上去做簽名。

為什麼要這麼麻煩？這是為了確保你的私鑰永遠不接觸網絡，安全至上。因為你永遠不知道你的 online 電腦是否中了木馬，會不會有人正在竊取你的數據。

這是去中心化面臨最困難的地方：你要自己保證安全，要為你的所有行為負責。

首先我們先看第一步，生成未簽名交易：

```
from pycoin.services.chain_so import ChainSoProvider
from pycoin.tx.tx_utils import create_tx

def genTx(fromAddress, toAddress, netcode):
    b = ChainSoProvider(netcode)
    spendables = b.spendables_for_address(fromAddress)
    if not spendables:
        print("{} no spenables tx.".format(fromAddress))
        return
    tx = create_tx(spendables, [toAddress])
    return tx.as_hex(include_unspents=True)

if __name__ == '__main__':
    fromAddress = 'mv3euCeE4zA5T4LVo1iMx3z29RGTDWG25U'
    toAddress = 'mqFKSAxHrSw21sXTbrFoXo7eMZfw6YHTaZ'
    netcode = 'XTN'

    rawTx = genTx(fromAddress, toAddress, netcode)
    print('Tx:', rawTx)
```

由於我們要測試轉帳，所以這次我們不用 BTC 主網，而是使用測試網 XTN，你可以去 [faucet](https://coinfaucet.eu/en/btc-testnet/) 免費申請測試幣。

同樣在 Python3 環境下執行，注意由於你的測試地址可能不同、或者該地址的餘額已經被轉走，你可能得到不一樣的結果：

```
Tx: 01000000012774c3f5db907b9d296771cfa05b3adbfcf14521666fec70e27e81f570f1e5b90000000000ffffffff01a093e400000000001976a9146abc3fcc31799a5fb7f9e897afa471e502f3c6b788ac00000000b0bae400000000001976a9149f6022d0fcfe1c66b29fdfb034caf2b84a66e4dd88ac
```

把 01000000 開頭的這串數字複製到 U 盤，下篇告訴你怎麼用私鑰對這個交易進行簽名和廣播。

