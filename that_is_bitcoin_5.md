# Python - 生成自己的比特幣私鑰

比特幣發展到現在超過十年，它有趣的地方在於，所有的都東西都由你自己控制。

當然前提是你會編程，下面告訴你用 Python 來生成比特幣私鑰，這樣你便可以在離線的電腦上，確保安全無虞的來生成自己的私鑰。

注意在使用私鑰之前，要確保做好備份。

Python 有個 Library - [pycoin](https://github.com/richardkiss/pycoin)，雖然最新還只是 0.9a 版本，但因為是開源程序，所以大部分情況下足夠使用了，不過今天我們使用穩定的 0.8 版本。

在生成私鑰之前，我們還需要知道比特幣的 BIP39 協議，讓我們可以使用助記詞來作為種子，所以我們還會用到 Library - [Mnemonic](https://github.com/trezor/python-mnemonic)，用於生成助記詞。

我們先安裝 Mnemonic：

```
pip install mnemonic
```

然後就可以使用 mnemonic 了：

```
from mnemonic import Mnemonic

def genWords(sec=256):
    m = Mnemonic("english")
    code = m.generate(sec)
    return code

if __name__ == '__main__':
    words = genWords(256)
    print("Seed:", words)
```

是不是非常簡單，genWords 函數裡面只需要幾行代碼就可以生成助記詞，256 Bits 對應 24  個詞，你還可以選擇其他強度 128 Bits= 12 詞、160 Bits=15 詞、192 Bits=18 詞、224 Bits=21 詞。

我們再安裝 pycoin:

```
pip install pycoin
```

然後就可以開始使用 pycoin 了，在上面代碼的基礎上加上 genRoot 函數：

```
from mnemonic import Mnemonic
from pycoin.key.BIP32Node import BIP32Node

def genWords(sec=256):
    m = Mnemonic("english")
    code = m.generate(sec)
    return code

def genRoot(words, passphrase, netcode, private=False):
    seed = Mnemonic.to_seed(words, passphrase)
    root = BIP32Node.from_master_secret(seed)
    root._netcode = netcode
    rootkey = root.wallet_key(as_private=private)
    return rootkey

if __name__ == '__main__':
    words = genWords(256)
    print("Seed:", words)

    passphrase = ''
    netcode = 'BTC'
    master = genRoot(words, passphrase, netcode, True)
    print("Rootkey:", master)
```

同樣也只要數行代碼就可以得到主私鑰，以上代碼在 Python3 環境下執行，可以得到這樣的結果：

```
Seed: online dilemma food relax tip visual dismiss input purpose flip pilot gold render invite stable cluster retreat earn sleep coach worth senior umbrella mad
Rootkey: xprv9s21ZrQH143K3yjczqoyXUUuBKsbEJwnMXPYYMJLdqRf1ZHzDscR6pfk9S1HCZRDYyPF81rVGxvqHh4nwqN4RPA4s2qQthPxQ5SoukqVwjt
```

備份好種子或者私鑰，下篇告訴你怎麼使用主私鑰生成帳戶公鑰。

如果你觉得这篇文章对你有帮助，欢迎你给我捐赠比特币！  
![btc-qrcode](https://jkeu374190052.files.wordpress.com/2019/01/1546697811.png)  
3GcQRzfZ6pWwntYkJBBJsmVcENn7ZoM8Kt
