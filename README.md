# `eosc` EOSIO command-line swiss-army-knife

## Description
[点击查看中文](./README-cn.md)

`eosc` is a cross-platform (Windows, Mac and Linux) command-line tool
for interacting with an EOS.IO blockchain.

## Features

* Superset of `cleos` functionalities
* System contract interactions (EOS Mainnet)
* Capacity to craft any transaction
* Multisig facilities with added bells and whistles
* Integrated secure Vault, to sign transactions and broadcast them
* Supports offline signature flows, and cold wallets.
* Based on the `eos-go` library, and is easy to extend.

## Installation

1. Install from https://github.com/eoscanada/eosc/releases

**or**

2. Build from source with:

```bash
go get -u -v github.com/eoscanada/eosc/eosc
```

**or**

3. If you are on **MacOS** and use **Homebrew**:

```
brew install eoscanada/tap/eosc
```


## Getting started

Once installed run:

```
eosc vault create --import
```

to import your keys and follow instructions.

Then set your environment variable to the API URL of your choice, optionally setting some HTTP headers:

```
export EOSC_GLOBAL_API_URL=https://your-favorite-endpoint

export EOSC_GLOBAL_HTTP_HEADER_0="Authorization: bearer abcdef12323453452565676589"
export EOSC_GLOBAL_HTTP_HEADER_1="Origin: https://something...
```

Then you can run commands on the chain, ex:

```
eosc get info
eosc transfer fromaccnt toaccnt 0.0001 --memo "Sent with eosc"
```

## Environment variables

These are supported environment variables:

* All global flags (those you get from eosc –help) can be set with the following pattern: EOSC_GLOBAL_FLAG_NAME. The most useful are:
  * `EOSC_GLOBAL_WALLET_URL` -> `--wallet-url`
  * `EOSC_GLOBAL_VAULT_FILE` -> `--vault-file`

* All (sub)command flags map to the following pattern: EOSC_COMMAND_SUBCOMMAND_CMD_FLAG_NAME (ex: `EOSC_FORUM_POST_CMD_REPLY_TO` -> `eosc forum post --reply-to=...`

Special cases:
* `EOSC_GLOBAL_INSECURE_VAULT_PASSPHRASE` allows you to input the passphrase directly in an environment variable (useful for test automation, risky for most other uses)
* `EOSC_GLOBAL_HTTP_HEADER_0` (available for indexes 0 to 25)



## Documentation

* [Voting offline](./OFFLINE_VOTING.md)
* [Offline transaction signature & offline multisignature](./OFFLINE_TRANSACTION_SIGNATURE.md)

## Cryptographic primitives used

The cryptography used is NaCl
([C implementation](https://tweetnacl.cr.yp.to/), [Javascript port](https://github.com/dchest/tweetnacl-js),
[Go version, which we use](https://godoc.org/golang.org/x/crypto/nacl/secretbox)). And
key derivation is [Argon2](https://en.wikipedia.org/wiki/Argon2),
using the [Go library
here](https://godoc.org/golang.org/x/crypto/argon2).

You can inspect the crypto code in our codebase regarding the
`passphrase` implementation: [it is 61 lines](./vault/passphrase.go),
including blanks and comments.



## FAQ

Q: Why not use `cleos` instead ?

A: `cleos` is hard to compile, being in C++, as it requires a huge
toolchain.  `eosc` works on Windows where `cleos` doesn't.  `eosc`
contains a wallet inside, and is able to use it to sign some
transactions, `cleos` interfaces with yet another program (`keosd`) in
order to sign transactions, making it more complex to use. `eosc`
brings `keosd` and `cleos` together in a swiss-army-knify package.
