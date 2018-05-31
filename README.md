`eosc` command-line swiss-knife
-------------------------------

`eosc` is a cross-platform (Windows, Mac and Linux) command-line tool
for interacting with an EOS.IO blockchain.

It contains tools for voting and a Vault to securely store private
keys.

It is based on the `eos-go` library, used in `eos-bios`, the tool used
to launch the EOS mainnet.

This first release holds simple tools, but a whole `cleos`-like
swiss-knife is being developed and will be released shortly after
mainnet launch.  Source code for most operations is already available
in this repository.


eosc vault
----------

The `eosc vault` is a simple yet powerful EOS wallet.


New wallet
==========

Create a new wallet with:

```
$ eosc vault create --keys 2
Created 2 keys. Let's secure them before showing the public keys.
Enter passphrase to encrypt your keys:
Confirm passphrase:
Wallet file "./eosc-vault.json" created. Here are your public keys:
- EOS7MEGq9FVb2Ve4bsZEan1t146TKCyo8dKtLvihrNhGbPLCPLjXd
- EOS5QoyZwJvpPjmZAa3HdgTn2FdNABBffXLD95WPagiARmaAHMhin
```

If you are still within the ERC-20 periods, you could register such
keys on Ethereum, to be ported to mainnet.


Importing keys
==============

Or import an existing private key.  This key will be encrypted to disk with your password in the default `vault-wallet.json` file:

```
$ eosc vault import
Vault file not found, creating a new wallet
Type your first private key:        [insert your private key here]
Type your next private key or hit ENTER if you are done:
Keys imported. Let's secure them before showing the public keys.
Enter passphrase to encrypt your keys:
Confirm passphrase:
Wallet file "./eosc-vault.json" written. Here are your public keys:
- EOS5tb61aZMAfQqKDsnkscFq76JXxNdi7ZhkUmkVZUkU4zPzfeAFx
```

The vault operations do zero network calls and can be done offline.
The file is encrypted, can safely be archived and is future proof.

Cryptography used
=================

The cryptography used is NaCl
([Javascript](https://github.com/dchest/tweetnacl-js),
[Go version which we use](https://godoc.org/golang.org/x/crypto/nacl/secretbox)). And
key derivation is [Argon2](https://en.wikipedia.org/wiki/Argon2),
using the [https://godoc.org/golang.org/x/crypto/argon2](Go library
here).

You can inspect the crypto code in our codebase,
[it is 79 lines](./vault/passphrase.go), including blanks and
comments.


eosc vote
---------

With a vault defined locally, you can vote:

```
$ eosc vote producers youraccount eoscanadacom someotherfavo riteproducer -u http://mainnet.eoscanada.com
Enter passphrase to unlock vault:
Voter [youraccount] voting for: [eoscanadacom]
Done
```

This will sign your vote transaction locally, and submit the transaction to the network through the `http://mainnet.eoscanada.com` endpoint.  You can also point to some other endpoints that are on the main network.

NOTE: This will not work before mainnet launches!


Features in the work
--------------------

* Have `serve` default to asking a confirmation on the command line before signing each transactions.
* Add `--accept` to auto-accepts signatures (only if listening on a local address).
* Shamir Secret Sharing-encoded wallets, with support for fast and easy multi-sig, even through a secure network.
* A full suite of tools to help Block Producers, developers, end users, in their every day life.


FAQ
---

Q: Why not use `cleos` instead ?

A: `cleos` is hard to compile, being in C++, as it requires a huge
toolchain.  It uses sha512 key derivation, which is faster to brute
force, the Argon2 key derivation is stronger and would take an
attacker a *lot* more efforts to bruteforce. `eosc` works on Windows
where `cleos` doesn't.  `eosc` contains a wallet inside, and is able
to use it to sign some transactions, `cleos` interfaces with yet
another program (`keosd`) in order to sign transactions, making it
more complex to use. `eosc` brings `keosd` and `cleos` together in a
swiss-knify package.
