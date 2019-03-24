# Python - 生成比特幣地址

上篇我們擁有了帳戶公鑰，現在看看如何用公鑰生成比特幣地址。

比特幣地址就是用來收幣的地址，通常有 Legacy 地址由 1 開頭，或者隔離見證 Segwit 地址由 3 開頭。

## 我們先來看看 Legacy 地址的生成：

```
from pycoin.key.key_from_text import key_from_text

def genAddressLegacy(pubkey, path, netcode):
    master = key_from_text(pubkey)
    key = master.subkey_for_path(path)
    key._netcode = netcode
    address = key.address()
    return address

if __name__ == '__main__':
    account = 'xpub6CcZXHumpbkN75Vb8r4yFfUSQCREy6PozhoVKvqCF7ZyyDkg2apkW4KzSNikAFuWnVtRVqFLyvWqfuYGv5nryDs8HVN1HCyBb1cKzfVqSjW'
    netcode = 'BTC'
    index = 0
    for i in range(index, index + 10):
        path = '0/' + str(i)
        address = genAddressLegacy(account, path, netcode)
        print(path, ":", address)
```

跟生成公鑰的代碼很類似，區別只在於 key.address()，上述程序在 Python3 環境下執行，會得到 10 個地址：

```
0/0 : 1A9bAsQ4RG7rEwpMCy2cfD1AFwKpfAu8Tm
0/1 : 1ANQHqLRbSbvnJiASfUPAahEaq1gT5zbKh
0/2 : 1LTyYK1VF76b91bV4iwXXqwjxVx7NpE1pq
0/3 : 1MjpkBCCNunmKgM8ESYNcVNth4pPUnxAC4
0/4 : 1Nxgc6c42ahYUY3kcEWLFZyzV9PnianAyj
0/5 : 1LXmvq75gtC6iGm6A7APEZUTFDmmhTADPQ
0/6 : 16qXUtKe5u2uRiCBBfsYTkPLFzDADyHyGx
0/7 : 1AJRE398kLqHZ2SqK4cYY5Z34WZ1DNENxr
0/8 : 1FrCiBL8BmgqqqS2H67LAr62cXEmTEVGVs
0/9 : 19SbvFiPUKjnGXpqHc3f177bgxNRyP6tcQ
```

## 我們再來看看隔離見證 Segwit 地址的生成：

隔離見證的地址稍微麻煩一點，它要用腳本方式 P2SH 付款，所以我們要先對公鑰進行 Hash160，然後構造付給這個公鑰的 Script。

```
    hash160_c = key.hash160(use_uncompressed=False)
    script = ScriptPayToAddressWit(b'\0', hash160_c).script()
    address = address_for_pay_to_script(script, key._netcode)
```

我們把上面這幾行代碼替換掉上面 Legacy 版本中取地址的部分，看看修改後的完整代碼：

```
from pycoin.key.key_from_text import key_from_text
from pycoin.tx.pay_to.ScriptPayToAddressWit import ScriptPayToAddressWit
from pycoin.ui import address_for_pay_to_script

def genAddressSegwit(pubkey, path, netcode):
    master = key_from_text(pubkey)
    key = master.subkey_for_path(path)
    key._netcode = netcode

    hash160_c = key.hash160(use_uncompressed=False)
    script = ScriptPayToAddressWit(b'\0', hash160_c).script()
    address = address_for_pay_to_script(script, key._netcode)

    return address

if __name__ == '__main__':
    account = 'xpub6CcZXHumpbkN75Vb8r4yFfUSQCREy6PozhoVKvqCF7ZyyDkg2apkW4KzSNikAFuWnVtRVqFLyvWqfuYGv5nryDs8HVN1HCyBb1cKzfVqSjW'
    netcode = 'BTC'
    index = 0
    for i in range(index, index + 10):
        path = '0/' + str(i)
        address = genAddressSegwit(account, path, netcode)
        print(path, ":", address)
```

我們在 Python3 環境下運行看看效果：

```
0/0 : 356Jm8bmryesDE19XMc23JfT2nJJehE7My
0/1 : 3AFAWyw18jJeb8MgJxsVdniUMcMtm6RF4B
0/2 : 39VtDwsAeunEE1oRqfhZ9P38p96rQxeR7g
0/3 : 3FkhtEPbYTDkDsowo7fYqo34VpmmDRiaqm
0/4 : 39jHpS5uctGJWxYg2E4dPEgP1kzjyX565f
0/5 : 32nVMtvT1MVXv6AZtQtKgKPoZHx1dnQVK4
0/6 : 3KydboYijXmQi3ndem9pseMp5Z1sRUfSbA
0/7 : 3HVCXzq5wKtsFTEcwihHD2MMWKPNzQwgVu
0/8 : 3PDCEEx1SYQoWpePrDMVwfUyouQ4mWjS7f
0/9 : 3DpVU46iyTypgEkFrTrTzRByqwmPzTJZXC
```

基於這個方式，你可以生成幾乎無數個地址。
