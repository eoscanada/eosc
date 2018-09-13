---
title: Global Flags
linktitle: Global Flags
description: Hugo's CLI is fully featured but simple to use, even for those who have very limited experience working from the command line.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [getting started]
keywords: [usage,livereload,command line,flags]
menu:
  docs:
    parent: "global-settings"
    weight: 40
weight: 40
sections_weight: 40
draft: false
aliases: [/overview/usage/,/extras/livereload/,/doc/usage/,/usage/]
toc: true
---

Flags:

  -u, --api-url string              API endpoint of eos.io blockchain node (default "https://mainnet.eoscanada.com")
  
      --delay-sec int               Set time to wait before transaction is executed, in seconds. Defaults to 0 second.
      --expiration int              Set time before transaction expires, in seconds. Defaults to 30 seconds. (default 30)
  -h, --help                        help for eosc
      --kms-gcp-keypath string      Path to the cryptoKeys within a keyRing on GCP
      --offline-chain-id string     Chain ID to sign transaction with. Use all --offline- options to sign transactions offline.
      --offline-head-block string   Provide a recent block ID (long-form hex) for TaPoS. Use all --offline options to sign transactions offline.
      --offline-sign-key strings    Public key to use to sign transaction. Must be in your vault or wallet. Use all --offline- options to sign transactions offline.
  -p, --permission strings          Permission to sign transactions with. Optionally specify more than one, or separate by comma
      --skip-sign                   Do not sign the transaction. Use with --write-transaction.
      --sudo-wrap                   Wrap the transaction in a eosio.sudo exec. Useful to BPs, with --write-transaction and --skip-sign to then submit as a multisig proposition.
      --vault-file string           Wallet file that contains encrypted key material (default "./eosc-vault.json")
      --wallet-url strings          Base URL to wallet endpoint. You can pass this multiple times to use the multi-signer (will use each wallet to sign multi-sig transactions).
      --write-transaction string    Do not broadcast the transaction produced, but write it in json to the given filename instead.
