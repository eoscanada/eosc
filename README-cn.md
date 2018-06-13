`eosc` 命令行工具 (瑞士军刀)
----------------------------

`eosc` 是一个跨平台的命令行工具 (支持 Windows, Mac 和 Linux)，你可以通过这个工具实现与 EOS.IO 区块链的交互。

这个工具可以用来投票，同时也可以作为存储私钥的钱包来使用。

本工具依赖 `eos-go` 库，而用来启动 EOS 主网的 `eos-bios` 使用的也是这个库。

当前版本只包含了最简单的命令，而类似 `cleos` 完整的功能还正在开发中，我们将会在主网上线后发布新的版本。几乎所有的执行源码都可以在当前仓库中找到。


安装
----

1. 安装包下载地址 https://github.com/eoscanada/eosc/releases

**或**

2. 基于源码进行构建:

    go get -u -v github.com/eoscanada/eosc/eosc


投票!
-----

安装完成后运行如下命令导入你的私钥:

```
eosc vault create --import
```

接着运行如下命令获取投票的帮助信息:

```
eosc vote --help
```

运行如下命令为你的候选者投票:

```
eosc vote producers [your account] [producer1] [producer2] [producer3]
```

确保你的版本为 `v0.7.0` 或高于此版本:

```
eosc version
```

阅读以下内容获取更详细的信息。


eosc 私钥钱包
-------------

`eosc vault` 是一个简单但却非常强大的 EOS 钱包。



将私钥导入创建的钱包
====================

你可以通过 `vault create --import` 命令 **导入外部生成的私钥**。

该交互式命令会提示粘贴你的私钥，导入时私钥并不会 **显示** 在屏幕上,导入完成后命令将会回显公钥（由导入的私钥衍生出来）以方便你进行交叉验证。

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

你的私钥将采用口令 (passphrase) 进行加密保护，参考后面的章节获取该工具所使用的加密库。



创建钱包的同时生成私钥
======================

创建钱包的同时生成 2 个私钥:

```
$ eosc vault create --keys 2
Created 2 keys. Let's secure them before showing the public keys.
Enter passphrase to encrypt your keys:
Confirm passphrase:
Wallet file "./eosc-vault.json" created. Here are your public keys:
- EOS7MEGq9FVb2Ve4bsZEan1t146TKCyo8dKtLvihrNhGbPLCPLjXd
- EOS5QoyZwJvpPjmZAa3HdgTn2FdNABBffXLD95WPagiARmaAHMhin
```


将私钥添加至已存在的钱包
========================

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

你可以在没有网络或离线的环境下运行钱包命令。
生成的钱包文件 `eosc-vault.json` 是被加密处理的，因此你可以放心的对它进行存储。

投票
----

有了 `eosc-vault.json` 你就可以投票了:

```
$ eosc vote producers youraccount eoscanadacom someotherfavo riteproducer
Enter passphrase to unlock vault:
Voter [youraccount] voting for: [eoscanadacom]
Done
```

上面的命令会在本地对投票交易进行签名处理，然后通过节点 `https://mainnet.eoscanada.com` 将交易发送至主网络。
你也可以通过 `-u` 或 `--api-url` 参数指定主网上的其他节点服务器。

在这里你可以查看你的账户信息 https://eosauthority.com/account



使用的加密库
------------

加密算法采用的是 NaCl
([C 语言实现](https://tweetnacl.cr.yp.to/), [Javascript 实现](https://github.com/dchest/tweetnacl-js),
[Go 版本, 本工具所依赖的库](https://godoc.org/golang.org/x/crypto/nacl/secretbox))。
和密钥衍生工具 [Argon2](https://en.wikipedia.org/wiki/Argon2),
采用的是 [Go 语言实现](https://godoc.org/golang.org/x/crypto/argon2)。

你可以在我们的代码库中审查与 `passphrase` 口令实现相关的代码: [第 61 行](./vault/passphrase.go)，包括空行及注释。




正在开发的新功能
----------------

* 通过 `serve` 选项使得命令行工具可以在签名每一笔交易前提示确认信息。
* 增加 `--accept` 参数以支持自动接收签名 (仅限监听本地地址)
* 基于 [Shamir Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing) 编码的钱包，支持快速和方便的多重签名。
* 一个完整的工具包，用来帮助区块生产者、开发人员以及终端用户完成他们的日常工作。


问答
----

问题: 为什么不使用 `cleos` ?

答案：`cleos` 使用 C++ 编写，由于依赖太多的工具链而很难被编译。`eosc` 可以在 Windows 上使用，而 `cleos` 却不可以。`eosc` 包含了一个内部钱包，可以很方便的用来签名交易，但 `cleos` 命令则需要借助 (`keosd`) 才可以实现对交易的签名，因此它很难使用。而 `eosc` 将 `cleos` 和 `keosd` 这两个工具整合为一个方便使用的命令行工具。

ERRATA: Previouly, you could read `It uses sha512 key derivation,
which is faster to brute force, the Argon2 key derivation is stronger
and would take an attacker a *lot* more efforts to bruteforce.`. It
was incorrect as `cleos` generates a big random password for you which
is effectively very hard to brute force, no matter which derivation
algo you are using.  The fact that `cleos` doesn't allow you to choose
your passphrase is a difference, but mostly in usability. You need to
store that large password somewhere, right?
