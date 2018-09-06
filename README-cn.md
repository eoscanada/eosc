`eosc` 像瑞士军刀一样的多功能命令行工具
----------------------------------

`eosc` 是一个跨平台的命令行工具 (支持 Windows, Mac 和 Linux)，你可以通过这个工具实现与 EOS.IO 区块链的交互。

这个工具可以用来投票，同时也可以作为存储私钥的钱包来使用。

本工具依赖 `eos-go` 库，`eos-bios` 使用的也是这个库。

当前版本只包含了最简单的命令，而类似 `cleos` 完整的功能还正在开发中，我们将会在主网上线后发布新的版本。几乎所有的执行源码都可以在当前仓库中找到。


安装
----

1. 安装包下载地址 https://github.com/eoscanada/eosc/releases

**或**

2. 基于源码进行构建:

```bash
go get -u -v github.com/eoscanada/eosc/eosc
```


投票!
-----

安装完成后运行下面命令来导入你的私钥:

```
eosc vault create --import
```

接着再运行下面命令获取投票的帮助信息:

```
eosc vote --help
```

然后运行如下命令为你的候选者投票: 

```
eosc vote producers [your account] [producer1] [producer2] [producer3]
```
命令中参数依次为: [你的账户名] [你要投票的BP1] [你要投票的BP2] [你要投票的BP3] （你共可以为30个BP投票）

确保你的版本是 `v0.7.0` 以上的，你可以用下面的命令检查你的版本:

```
eosc version
```

阅读以下内容获取更详细的信息。


eosc 私钥钱包
-------------

`eosc vault` 是一个简单但却非常强大的 EOS 钱包。



将已有私钥导入新建的钱包
====================

你可以通过 `vault create --import` 命令 **导入外部生成的私钥**。

这个交互式命令会要你粘贴你的私钥，导入时私钥并 **不会显示在屏幕上**，导入完成后命令将会再次显示（由导入的私钥衍生出来）公钥，以方便你进行交叉验证。

```
$ eosc vault create --import -c "Imported key"
Enter passphrase to encrypt your vault: ****
Confirm passphrase: ****
Type your first private key: ****
Type your next private key or hit ENTER if you are done:
Imported 1 keys. Let's secure them before showing the public keys.
Wallet file "./eosc-vault.json" created. Here are your public keys:
- EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV
```

你的私钥将采用密语 (passphrase) 进行加密保护，参考后面的章节获取该工具所使用的加密库。



创建钱包的同时生成公钥
===================

创建钱包的同时生成 2 个公钥:

```
$ eosc vault create --keys 2
Created 2 keys. Let's secure them before showing the public keys.
Enter passphrase to encrypt your keys:
Confirm passphrase:
Wallet file "./eosc-vault.json" created. Here are your public keys:
- EOS7MEGq9FVb2Ve4bsZEan1t146TKCyo8dKtLvihrNhGbPLCPLjXd
- EOS5QoyZwJvpPjmZAa3HdgTn2FdNABBffXLD95WPagiARmaAHMhin
```


将已有的私钥添加至已存在的钱包
=========================

将外部生成的私钥导入当前钱包:

```
$ eosc vault add
Vault file not found, creating a new wallet
Type your first private key:        [insert your private key here, NOT shown]
Type your next private key or hit ENTER if you are done:
Keys imported. Let's secure them before showing the public keys.
Enter passphrase to encrypt your keys:
Confirm passphrase:
Wallet file "./eosc-vault.json" written. Here are your ADDED public keys:
- EOS5tb61aZMAfQqKDsnkscFq76JXxNdi7ZhkUmkVZUkU4zPzfeAFx
Total keys writteN: 3
```

你可以离线运行钱包命令。
生成的钱包文件 `eosc-vault.json` 是被加密处理的，因此你可以放心的把它存起来，他以后就是代替你私钥的凭证。

投票
----

`eosc-vault.json` 成功生成之后，你就可以投票了:

```
$ eosc vote producers youraccount eoscanadacom someotherfavo riteproducer
Enter passphrase to unlock vault:
Voter [youraccount] voting for: [eoscanadacom]
Done
```

上面的命令会在本地对投票交易进行签名处理，然后通过节点 `https://mainnet.eoscanada.com` 将交易发送至主网络。
你也可以通过 `-u` 或 `--api-url` 参数指定主网上的其他节点服务器。

在这里你可以查看你的账户信息 https://eosquery.com



使用的加密库
----------

加密算法采用的是 NaCl
([C 语言实现](https://tweetnacl.cr.yp.to/), [Javascript 实现](https://github.com/dchest/tweetnacl-js),
[Go 版本, 本工具所依赖的库](https://godoc.org/golang.org/x/crypto/nacl/secretbox))。
和密钥衍生工具 [Argon2](https://en.wikipedia.org/wiki/Argon2),
采用的是 [Go 语言实现](https://godoc.org/golang.org/x/crypto/argon2)。

你可以在我们的代码库中审查与 `passphrase` 口令实现相关的代码: [第 61 行](./vault/passphrase.go)，包括空行及注释。




## 离线交易签名


`eosc` 具备管理一个 airtight 脱机钱包所需的一切。
这是一个示例用法。 在 `/tmp/airtight` 中, 运行:

* `eosc vault create --import --comment "Airtight wallet"`

它会输出:

```
PLEASE READ:
We are now going to ask you to paste your private keys, one at a time.
They will not be shown on screen.
Please verify that the public keys printed on screen correspond to what you have noted

Paste your first private key:       <em>[PASTE PRIVATE KEY HERE]</em>
- Scanned private key corresponding to EOS5J53bNtH5H1bJTyP237fr5LF6eQSQohfGY5iMCgpC4HpXApJBr
Paste your next private key or hit ENTER if you are done:           <em>[PASTE SECOND PRIVATE KEY HERE]</em>
- Scanned private key corresponding to EOS6tgsdv6S7N1GYWgX8QEBAsanAwXwuaEkv11GGtteyk5ELqSzVP
Paste your next private key or hit ENTER if you are done:           <em>[HIT ENTER HERE]</em>
Imported 2 keys.

You will be asked to provide a passphrase to secure your newly created vault.
Make sure you make it long and strong.

Enter passphrase to encrypt your vault:       <em>[ENTER YOUR PASSPHRASE HERE]</em>
Confirm passphrase:                           <em>[RE-ENTER YOUR PASSPHRASE HERE]</em>

Wallet file "./eosc-vault.json" written to disk.
Here are the keys that were ADDED during this operation (use `list` to see them all):
- EOS5J53bNtH5H1bJTyP237fr5LF6eQSQohfGY5iMCgpC4HpXApJBr
- EOS6tgsdv6S7N1GYWgX8QEBAsanAwXwuaEkv11GGtteyk5ELqSzVP
Total keys stored: 2
```

这一操作会创建 `eosc-vault.json` 钱包，而且完全可以离线执行，不需要碰到网络上的任何东西。

好，那让我们来实际操作它试试。从一个有互联网访问权限的计算机发出这些命令（就是说这些命令无法访问你的 eosc 保险库钱包）:

* `eosc transfer airtight11111 cancancan123 1.0000 -m "Can't say I haven't paid you now" --write-transaction transaction.json --skip-sign --expiration 3600`

输出就会是这样:

```
{
  "expiration": "2018-08-27T04:18:43",
  "ref_block_num": 43579,
  "ref_block_prefix": 3787183056,
  "max_net_usage_words": 0,
  "max_cpu_usage_ms": 0,
  "delay_sec": 0,
  "context_free_actions": [],
  "actions": [
    {
      "account": "eosio.token",
      "name": "transfer",
      "authorization": [
        {
          "actor": "airtight1111",
          "permission": "active"
        }
      ],
      "data": {
        "from": "airtight1111",
        "to": "cancancan123",
        "quantity": "0.0001 EOS",
        "memo": "Can't say I haven't paid you now"
      }
    }
  ],
  "transaction_extensions": [],
  "signatures": [],
  "context_free_data": []
}
---
Transaction written to "../transaction.json"
Sign offline with: --offline-chain-id=5fff1dae8dc8e2fc4d5b23b2c7665c97f9e9d8edf2b6485a86ba311c25639191
Above is a pretty-printed representation of the outputted file
```

好，现在复制 `transaction.json` 文件到另外一台计算机上。注意上面输出的 `--offline-chain-id=...`。

在一个 airtight 计算机上运行:

* `eosc tx sign path/to/transaction.json --offline-sign-key EOS6tgsdv6S7N1GYWgX8QEBAsanAwXwuaEkv11GGtteyk5ELqSzVP  --offline-chain-id=5fff1dae8dc8e2fc4d5b23b2c7665c97f9e9d8edf2b6485a86ba311c25639191 --write-transaction signedtx.json`

你想输入几个 `--offline-sign-key` 就可以输几个，你可以用 `,` (逗号) 来分隔每个条目，例如: `--offline-sign-key EOS6tgsdv6S7N1GYWgX8QEBAsanAwXwuaEkv11GGtteyk5ELqSzVP,EOS5J53bNtH5H1bJTyP237fr5LF6eQSQohfGY5iMCgpC4HpXApJBr`。这时的输出是:

```
Enter passphrase to decrypt your vault:
{
  "expiration": "2018-08-27T04:18:43",
  "ref_block_num": 43579,
  "ref_block_prefix": 3787183056,
  "max_net_usage_words": 0,
  "max_cpu_usage_ms": 0,
  "delay_sec": 0,
  "context_free_actions": [],
  "actions": [
    {
      "account": "eosio.token",
      "name": "transfer",
      "authorization": [
        {
          "actor": "airtight1111",
          "permission": "active"
        }
      ],
      "data": "104208b93197af33304498064d83a641010000000000000004454f53000000002043616e277420736179204920686176656e2774207061696420796f75206e6f77"
    }
  ],
  "transaction_extensions": [],
  "signatures": [
    "SIG_K1_KbsvdjfXTzbn7DaghCbFyaxXnv9BHJaAVWZkeJ3AbvnznNgipSLj7BbZajZ1ECZEWKSMFtMbuqfUWQGq2tDzW2n7Fz5KTV",
    "SIG_K1_K24N3kUfTLMJGwUUWV2qNiyhxAJsbVtpXh2KUj3SHGLTivPBru46QYX9v9gQj7G2yKHp6eZ786hHfJwuzjeZFF7atPfpTY"
  ],
  "context_free_data": []
}
---
Transaction written to "signedtx.json"
Sign offline with: --offline-chain-id=5fff1dae8dc8e2fc4d5b23b2c7665c97f9e9d8edf2b6485a86ba311c25639191
Above is a pretty-printed representation of the outputted file
```

你现在可一个把 `signedtx.json` 放到一个在线的计算机上了，然后运行:

* `eosc tx push signedtx.json`

这个应该是成功的，假如没有，请参照下面的解释:

* `transaction bears irrelevant signatures from these keys: [\"EOS5J53bNtH5H1bJTyP237fr5LF6eQSQohfGY5iMCgpC4HpXApJBr\"]`: 这意味着你无需使用此密钥签名即可授权你在交易的操作中指定的 `actor@permission`。

* `UnAuthorized` errors: 这是说你没有使用交易操作中指定授权
`actor @ permission` 所需的密钥对交易进行签名。

* `Transaction's reference block did not match.`: 这是说你的交易没有创建在你想推送的链上。

* `expired transaction`: 你需要在原始 `--expiration` 的时间范围内完成整个循环。 
如果你给 `--expiration` 超过一个小时，那请注意你只能在到期的最后一小时内将交易提交到链上。


* 一些其他 assertion（断言）错误，就是你尝试在线签名时发生的常见错误，请根据相关合约进行调查。



### 离线 multisig（多签名）


为增加离线交易签名的安全性，请考虑使用EOS区块链的 `multisig` 工具。

你可以在链上（通过`eosc multisig proposal`）提出一个交易，所有人都可以审查，
并准备一个离线交易来批准该交易。 使用这种方法，你可以拥有 X 个 airtight 的计算机，
每个计算机都提供拼图中的一块，由不同的人控制，无需缝合交易文件或将签名收集到一个地方。


常见问题
-------

问题: 为什么不使用 `cleos` ?

答案：`cleos` 使用 C++ 编写，由于依赖太多的工具链而很难被编译。`eosc` 
可以在 Windows 上使用，而 `cleos` 却不可以。`eosc` 包含了一个内部钱包，
可以很方便的用来签名交易，但 `cleos` 命令则需要借助 (`keosd`) 才可以实现对交易的签名，
因此它很难使用。而 `eosc` 将 `cleos` 和 `keosd` 这两个工具整合为一个方便使用的命令行工具。
