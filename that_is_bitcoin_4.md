# 这才是比特币（4）- 交易签名

## 交易数据格式

我们先来看看一个比特币区块链上的交易长什么样子

````
0100000001a33559f42f5be03a5466c3394dbf0f6c2ece5a53649289ef89bd2cdb4f962845010000006a4730440220729c8d2177924cab79b8aab33dc9b3a4d6d19cf59acc770b49913a689d7f1d0b02202a5ce3522007fdfebd33e35ffb94c6ca1f031b943086e80dc5c2881069bc319e0121025e1f7a423049dcfc56aea40923a9fd8dfe8e15c294434e76f5d424146ef3ad60ffffffff0210270000000000001976a914450c4e370fe7b9e06c4e48791c6c3a308ffc91bd88aca2670b00000000001976a9140c8da7439d4a04dafde5ddb47feb94e09b8dbf5888ac00000000
````

以上这串数据是根据下面的格式拼凑起来的，交易数据格式:

|Length |Name     
|-------|---
|4Bytes |Version
|1Bytes |TxInCount
|Var    |TxIn
|1Bytes |TxOutCount
|Var    |TxOut
|4Bytes |LockTime

TxIn:

|Length |Name
|-------|---
|Var    |UTXO.TxId
|4Bytes |UTXO.Index
|Var    |SignatureScript
|4Bytes |Sequence

TxOut:

|Length |Name
|-------|---
|8Bytes |Value
|Var    |PkScript

## 输入数据

一个交易的输入数据，必须是链上的未花费交易（UTXO），所以我们每次交易之前，都需要先确认有哪些 UTXO 可以使用。

我们可以从 Bitcoin Core 节点中查询 UTXO，也可以从其他开放的 API 中获取这个列表。例如使用 chain.so 的接口：

https://chain.so/api/v2/get_tx_unspent/BTCTEST/mgfL6c3g2vwn2KhE6nm3TSgxMyKsEGyzHt

可以得到未花费交易的 txid、交易序号 output_no 和解锁脚本 script_hex，把这 3 个数据组合起来成为一个 TxIn。

## 输出数据

每次交易的输出，除了目标地址和目标数量外，还需要指定找零地址和找零数量。这是因为比特币的 UTXO 机制决定了每个交易都必须把输入的 UTXO 都消费完，不消费部分就是交易手续费。

把目标数量和目标地址生成的解锁脚本 PkScript 组合起来成为一个 TxOut。

## 未签名数据

为了保护私钥安全，通常我们会选择离线签名，所以我们先在一台在线的电脑中制作未签名数据，注意在制作为签名数据的过程中，不需要私钥参与。

我们把输入数据和输出数据组合成起来，补充上版本号、时间、数量等信息，就得到一个完整的未签名交易，把这个文件格式化未字符串数据保存到未签名文件 unsign.txn 中。

## 离线签名

我们把在线电脑中的 unsign.txn 复制到 U 盘，再把 U 盘插入离线电脑进行签名。

签名过程分为四步：

1. 把当前输入数据中 UTXO 的解锁脚本保留，其他 UTXO 解锁脚本清零；
2. 对以上数据使用私钥进行签名，再拼上公钥，就得到签名数据=签名+公钥；
3. 用签名数据替换输入数据中 UTXO 的解锁脚本；
4. 重复以上三步直到所有 UTXO 签名完成。

## 广播交易

把离线电脑中签名文件 signed.txn 用 U 盘复制到在线电脑中进行广播。

广播可以使用 Bitcoin Core 节点，也可以使用其他开放的 API 进行广播。例如同样可以使用 chain.so 接口：

https://chain.so/api/v2/send_tx/BTCTEST

还可以使用 Electrum -> Tools -> Load Transaction -> From file 进行广播。

或者使用其他的网站进行广播 https://live.blockcypher.com/btc/pushtx

等待约 10 分钟直到交易被确认，恭喜你完成了一笔比特币转账交易。

完整代码请见 [Github](https://github.com/jkeu/that_is_bitcoin)

如果你觉得这篇文章对你有帮助，欢迎你给我捐赠比特币！  
![btc-qrcode](https://jkeu374190052.files.wordpress.com/2019/01/1546697811.png)  
3GcQRzfZ6pWwntYkJBBJsmVcENn7ZoM8Kt
