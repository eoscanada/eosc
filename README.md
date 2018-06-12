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


Installation
------------

1. Install from https://github.com/eoscanada/eosc/releases

**or**

2. Build from source with:

```bash
go get -u -v github.com/eoscanada/eosc/eosc
```



Vote now!
---------

Once installed run:

```
eosc vault create --import
```

to import your keys and follow instructions.

Then run:

```
eosc vote --help
```

and run somethig like this:

```
eosc vote producers [your account] [producer1] [producer2] [producer3]
```

Make sure you have version `v0.7.0` or higher:

```
eosc version
```

Read more below for details.


eosc vault
----------

The `eosc vault` is a simple yet powerful EOS wallet.



Import keys in new wallet
=========================

You can **import externally generated private keys** with `vault create --import` call.

You will then be asked to paste your private key in an interactive
prompt.  The private key is **not shown** on screen during the
import. It will display the corresponding public key which you should
cross-check.

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

Your private keys will be encrypted using your passphrase. See below
for cryptographic primitives used.



New wallet with new keys
========================

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


Add keys to an existing vault
=============================

To add an externally generated private key to an existing vault, use:

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

The vault operations do zero network calls and can be done offline.
The file is encrypted, can safely be archived and is future proof.

Casting your votes
------------------

With an `eosc-vault.json`, you can vote:

```
$ eosc vote producers youraccount eoscanadacom someotherfavo riteproducer
Enter passphrase to unlock vault:
Voter [youraccount] voting for: [eoscanadacom]
Done
```

This will sign your vote transaction locally, and submit the
transaction to the network through the `https://mainnet.eoscanada.com`
endpoint.  You can also point to some other endpoints that are on the
main network with `-u` or `--api-url`.

Find what your account it on https://eosquery.com/



Cryptographic primitives used
-----------------------------

The cryptography used is NaCl
([C implementation](https://tweetnacl.cr.yp.to/), [Javascript port](https://github.com/dchest/tweetnacl-js),
[Go version, which we use](https://godoc.org/golang.org/x/crypto/nacl/secretbox)). And
key derivation is [Argon2](https://en.wikipedia.org/wiki/Argon2),
using the [Go library
here](https://godoc.org/golang.org/x/crypto/argon2).

You can inspect the crypto code in our codebase regarding the
`passphrase` implementation: [it is 61 lines](./vault/passphrase.go),
including blanks and comments.




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
toolchain.  `eosc` works on Windows where `cleos` doesn't.  `eosc`
contains a wallet inside, and is able to use it to sign some
transactions, `cleos` interfaces with yet another program (`keosd`) in
order to sign transactions, making it more complex to use. `eosc`
brings `keosd` and `cleos` together in a swiss-knify package.

ERRATA: Previouly, you could read `It uses sha512 key derivation,
which is faster to brute force, the Argon2 key derivation is stronger
and would take an attacker a *lot* more efforts to bruteforce.`. It
was incorrect as `cleos` generates a big random password for you which
is effectively very hard to brute force, no matter which derivation
algo you are using.  The fact that `cleos` doesn't allow you to choose
your passphrase is a difference, but mostly in usability. You need to
store that large password somewhere, right?
