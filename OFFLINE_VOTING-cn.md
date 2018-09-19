如何使用 eosc 和密闭计算机进行投票
-----------------------------------------------

密闭计算机：Air Gapped Computer，
我们这里讲的离线签名方法需要两个计算机，其中一个不联网 （这个就叫密闭计算机），
通过把签名过程放在在密闭计算机上进行，我们就保证了密钥的安全。

本步骤分解是建立在假设你已经在主网上拥有一个帐号的前提下的
（我们接下来就称它为“YOURACCOUNT”），
在下面的步骤中你会用到与这个账号私钥相对应的公钥。
    
1. [在线计算机]从`eosc`发布页中下载的最新版本
   https://github.com/eoscanada/eosc/releases
   选择与您的密闭计算机操作系统兼容的版本：Windows，Mac或Linux
   （最好你能用占用空间小的Linux版）。

2. [在线计算机] 同时也给您的在线计算机下好`eosc`。
 
3. [密闭计算机] 把 `eosc` 拷到密闭计算机上
 
4. [密闭计算机] 创建新的密钥保险库，输入你的私钥:

 ```
eosc vault create --import --comment "Airtight wallet"
```

这一步会创建一个新的 `eosc-vault.json` 文件，
它里面加密了你的密钥，存在密闭计算机的磁盘上。

 5. [在线计算机]准备好你的`voteproducer`交易：
 
 ```
eosc vote producers YOURACCOUNT PRODUCER1 PRODUCER2 PRODUCER3... --write-transaction unsigned-tx.json --skip-sign --expiration 3600
```

这一步将创建一个包含 `voteproducer` 交易的 unsigned-tx.json 文件。
它不会试图对其签名（这将发生在密闭计算机上），而是设置一个小时的到期时间，
在这期间你需要完成文件的签名来回传输。

你可以使用 `eosc tx unpack unsigned-tx.json` 解压查阅其中的数据。

注意：你可以以同样的方式使用任何 `eosc` 命令来的输出这样待签名的交易（如
`eosc transfer`，`eosc system delegatebw`）。

6. [在线计算机]以任何方式将 `unsigned-tx.json` 文件复制到密闭计算机上。
    
7. [密闭计算机]通过用 Vault（保险库）中的与你私钥匹配的公钥（以及 YOURACCOUNT 帐户权限）代替 `PUBLIC_KEY` 来签署交易：

```
eosc tx sign unsigned-tx.json --offline-sign-key PUBLIC_KEY..... --write-transaction signed-tx.json
```

这个命令结束后将会产生 `signed-tx.json` 文件。

8. [密闭计算机]将 `signed-tx.json` 文件再放到在线计算机上。

9. [在线计算机]将签好了的交易推送到链上：

 ```
eosc tx push signed-tx.json
```

就这样，完事了！
 
### 故障排除

如果交易由于某种原因失败，你可以看看弹出的错误消息是否匹配以下之一：

 
* `UnAuthorized` errors：这意味着你没有使用交易 action 中指定授权 
  `actor@permission` 所需的密钥对交易成功进行签名。
  
* `expired transaction`：你需要在原本的 `--expiration` 的时间范围内完成整个签名过程。 
  如果你给 `--expiration` 设置超过一个小时，那你只能在到时间前的最后一小时内将交易提交到链上。
  
* `Transaction's reference block did not match.`：这是说交易的参考区块不匹配。
  你的交易可能没有创建在你想推送的链上。
  
* `transaction bears irrelevant signatures from these keys: [\"EOS5J5....\"]`：
  翻译：`交易中包含与这个密钥无关的签名：[\“EOS5J5 .... \”]`：
  意思是你不需要使用此密钥来签名授权在交易 action 中指定的 `actor @ permission`。
  
* 其他断言错误：如果你签名的时候是在线的，一般会出现这个错误，请调查相关合约。
