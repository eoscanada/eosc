---
title: Quick Start
linktitle: Quick Start
description: Install and run eosc
date: 2018-09-11
publishdate: 2018-09-11
categories: [getting started]
keywords: [quick start]
authors: ["EOS Canada"]
menu:
  docs:
    parent: "getting-started"
    identifier: quick-start-getting-started
    weight: 10
weight: 10
sections_weight: 10
draft: false
aliases: [/quickstart/,/overview/quickstart/]
toc: true
---

Use `eosc` to vote for your favorite Block Producers
------------

Once installed run (see : {{< ref "../install-eosc.md" >}} ):

```
eosc vault create --import
```

to import your keys and follow instructions.

Then run:

```
eosc vote --help
```

and run something like this:

```
eosc vote producers [your account] [producer1] [producer2] [producer3]
```

Using `eosc` as a cryptographic key vault (or wallet)
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

Find what your account is on https://eosq.app


