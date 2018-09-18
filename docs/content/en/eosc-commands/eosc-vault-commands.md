---
title: "`eosc vault` Commands"
linktitle: "`eosc vault`"
description:
date: 2018-09-05
publishdate: 2018-09-05
lastmod: 2018-09-05
keywords: []
menu:
  docs:
    parent: "eosc-commands"
    weight: 30
weight: 30
sections_weight: 30
draft: false
aliases: []
toc: true


---

```
Usage:
  eosc vault [command]

Available Commands:
  add         Add private keys to an existing vault taking input from the shell
  create      Create a new encrypted EOS keys vault
  export      Export private keys (and corresponding public keys) inside an eosc vault.
  list        List public keys inside an eosc vault.
  serve       Serves signing queries on a local port.

Flags:
  -h, --help   help for vault
```

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
