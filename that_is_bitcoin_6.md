# Python - 生成比特幣帳戶公鑰

比特幣地址格式經過多年的發展，基本上已經形成一定的共識。

首先是 BIP44 協議，因為分層確定性錢包 HD Wallet 的出現，也形成了對於錢包路徑使用上的規則：

```
m / purpose' / coin_type' / account' / change / address_index
```

purpose = 44，對於比特幣 coin_type = 0

所以通常主私鑰下的第一個地址是 m/44H/0H/0H/0/0

為了保證安全，通常我們生成一個 account 的公鑰，然後用這個公鑰就可以生成無數個地址，我們先用上篇生成的主私鑰來生成一個 account 的公鑰：

```
from pycoin.key.key_from_text import key_from_text

def genAccount(root, path, netcode):
    master = key_from_text(root)
    key = master.subkey_for_path(path)
    key._netcode = netcode
    pubkey = key.wallet_key()
    prvkey = key.wallet_key(as_private=True)
    return pubkey, prvkey

if __name__ == '__main__':
    master = 'xprv9s21ZrQH143K3yjczqoyXUUuBKsbEJwnMXPYYMJLdqRf1ZHzDscR6pfk9S1HCZRDYyPF81rVGxvqHh4nwqN4RPA4s2qQthPxQ5SoukqVwjt'
    path = '44H/0H/0H'
    netcode = 'BTC'
    pubkey, prvkey = genAccount(master, path, netcode)

    print("AccountPath:", path)
    print("AccountPrivate:", prvkey)
    print("AccountPublic:", pubkey)
```

在 Python3 環境下執行，可以得到：

```
AccountPath: 44H/0H/0H
AccountPrivate: xprv9ydD7nNszEC4tbR82pXxtXXhrAakZdfxdUstXYRagn316RRXV3WVxG1Wb649baZLSCrR7hff6VfF6GiQc2j3WsWP1eUJZX5CWo2HtCLxYTG
AccountPublic: xpub6CcZXHumpbkN75Vb8r4yFfUSQCREy6PozhoVKvqCF7ZyyDkg2apkW4KzSNikAFuWnVtRVqFLyvWqfuYGv5nryDs8HVN1HCyBb1cKzfVqSjW
```

現在你可以把帳戶公鑰拿到在線電腦上使用了，把它複製出來，下篇我們會用到這個公鑰來生成地址。

```
xpub6CcZXHumpbkN75Vb8r4yFfUSQCREy6PozhoVKvqCF7ZyyDkg2apkW4KzSNikAFuWnVtRVqFLyvWqfuYGv5nryDs8HVN1HCyBb1cKzfVqSjW
```


如果你觉得这篇文章对你有帮助，欢迎你给我捐赠比特币！  
![btc-qrcode](https://jkeu374190052.files.wordpress.com/2019/01/1546697811.png)  
3GcQRzfZ6pWwntYkJBBJsmVcENn7ZoM8Kt
