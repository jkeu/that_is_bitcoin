# 这才是比特币（3）- 生成钱包地址

比特币的去中心化和开源特性，让所有人都可以参与，并实现自己的逻辑。下面我们看看如何实现生成钱包地址，稍微搜索一下就可以找到很多描述地址生成过程的流程图，这里就不赘述了，直接看过程和代码，这里用 GO 语言做例子。

## 用椭圆曲线非对称加密算法生成私钥公钥对

````go
func newKeyPair(compress bool) (ecdsa.PrivateKey, []byte) {
  curve := btcec.S256()
  private, _ := ecdsa.GenerateKey(curve, rand.Reader)

  var pubKey []byte

  if !compress {
    format := byte(0x04)
    pubKey = append([]byte{format}, private.PublicKey.X.Bytes()...)
    pubKey = append(pubKey, private.PublicKey.Y.Bytes()...)
  } else {
    format := byte(0x02)
    if isOdd(private.PublicKey.Y) {
      format |= byte(0x01)
    }
    pubKey = append([]byte{format}, private.PublicKey.X.Bytes()...)
  }

  return *private, pubKey
}
````

早期的比特币用非压缩公钥65字节，现在由于区块数据越来越大，基本上都用压缩公钥33字节。

`非压缩公钥 = 0x04 + PublicKey.X + PublicKey.Y`

`压缩公钥 = 0x02(Y为偶数) or 0x03(Y为奇数) + PublicKey.X`

## 创建钱包

````go
type Wallet struct {
  PrivateKey ecdsa.PrivateKey
  PublicKey  []byte
  Compress   bool
}

func NewWallet() *Wallet {
  compress := true
  private, public := newKeyPair(compress)
  wallet := Wallet{private, public, compress}

  return &wallet
}
````

## 用公钥进行 SHA256 和 RipeMD160 运算得到公钥的 Hash

````go
func HashPubKey(pubKey []byte) []byte {
  publicSHA256 := sha256.Sum256(pubKey)

  RIPEMD160Hasher := ripemd160.New()
  _, err := RIPEMD160Hasher.Write(publicSHA256[:])
  if err != nil {
    log.Panic(err)
  }
  publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

  return publicRIPEMD160
}
````

`公钥 Hash = RipeMD160(SHA256(公钥))`

## 加上地址版本和公钥 Hash 的校验数据

`地址 = Base58(地址版本 + 公钥 Hash + 校验数据)`

### Legacy 版本

Legacy 地址版本用 0x00，生成的地址是 1 开头的。

````go
func (w Wallet) GetAddressPublicKey() []byte {
  pubKeyHash := HashPubKey(w.PublicKey)

  versionedPayload := append([]byte{byte(0x00)}, pubKeyHash...)
  checksum := checksum(versionedPayload)

  fullPayload := append(versionedPayload, checksum...)
  address := Base58Encode(fullPayload)

  return address
}
````

### Segwit-P2SH 版本

Segwit 地址版本用 0x05，生成的地址是 3 开头的。

````go
func (w Wallet) GetAddressScriptHash() []byte {
  pubKeyHash := HashPubKey(w.PublicKey)
  op := []byte{byte(0x00), byte(0x14)}
  pubKeyHash = append(op, pubKeyHash...)
  pubKeyHash = HashPubKey(pubKeyHash)

  versionedPayload := append([]byte{byte(0x05)}, pubKeyHash...)
  checksum := checksum(versionedPayload)

  fullPayload := append(versionedPayload, checksum...)
  address := Base58Encode(fullPayload)

  return address
}
````

### 校验数据

校验方法，对数据进行 2 次 SHA256，然后取前 4 个字节。

````go
func checksum(payload []byte) []byte {
  firstSHA := sha256.Sum256(payload)
  secondSHA := sha256.Sum256(firstSHA[:])

  return secondSHA[:addressChecksumLen]
}
````

### 地址使用 Base58 格式化输出

````go
func Base58Encode(b []byte) []byte {
  x := new(big.Int)
  x.SetBytes(b)

  answer := make([]byte, 0, len(b)*136/100)
  for x.Cmp(bigZero) > 0 {
    mod := new(big.Int)
    x.DivMod(x, bigRadix, mod)
    answer = append(answer, alphabet[mod.Int64()])
  }

  for _, i := range b {
    if i != 0 {
      break
    }
    answer = append(answer, alphabetIdx0)
  }

  alen := len(answer)
  for i := 0; i < alen/2; i++ {
    answer[i], answer[alen-1-i] = answer[alen-1-i], answer[i]
  }

  return answer
}
````

## 用私钥生成 WIF 格式显示

`WIF = Base58(0x08 + PrivateKey.D + 1 压缩位 + 4校验数据)`

非压缩 WIF 是 5 开头的，压缩 WIF 是 K 或 L 开头的。

````go
func (w Wallet) GetWif() []byte {
  versionedPayload := append([]byte{byte(0x80)}, w.PrivateKey.D.Bytes()...)
  if w.Compress {
    versionedPayload = append(versionedPayload, byte(0x01))
  }
  checksum := checksum(versionedPayload)

  fullPayload := append(versionedPayload, checksum...)
  address := Base58Encode(fullPayload)

  return address
}
````

## 在 main 开始调用上述流程

````go
func main() {
  log.Println("Generate Bitcoin address")

  myWallet := NewWallet()

  myAddressLegacy := myWallet.GetAddressPublicKey(false)
  log.Println("Address Legacy:", string(myAddressLegacy))

  myAddressSegWit := myWallet.GetAddressScriptHash(true)
  log.Println("Address Segwit:", string(myAddressSegWit))

  myWif := myWallet.GetWif()
  log.Println("Private:", string(myWif))
}
````

看看输出效果

````bash
2019/01/07 21:33:50 Generate Bitcoin address
2019/01/07 21:33:51 Address Legacy: 1ANkQNKgRHw3xthr2ctYf9mn9yjAfdNN9N
2019/01/07 21:33:51 Address Segwit: 33mWUz4jh54Jzs1nuE5qstvVagZuMvWBV4
2019/01/07 21:33:51 Private: L4vXVpQZT3zwv1GvwmscSZ3hkD4hEvrroM8oQ4SqggHJkAVAMdgc
````

你可以在 bitWallet 导入私钥 WIF，验证一下地址是不是跟我们生成的地址一样。

完整代码请见 [Github](https://github.com/jkeu/that_is_bitcoin)

如果你觉得这篇文章对你有帮助，欢迎你给我捐赠比特币！  
![btc-qrcode](https://jkeu374190052.files.wordpress.com/2019/01/1546697811.png)  
3GcQRzfZ6pWwntYkJBBJsmVcENn7ZoM8Kt
