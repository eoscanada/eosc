---
title: tx Commands Overview
linktitle: tx Commands Overview
description: Overview of the `tx` section of eosc
date: 2018-09-27
publishdate: 2018-09-27
lastmod: 2018-09-27
categories: [eosc-tx-commands]
keywords: []
menu:
  docs:
    parent: "eosc-tx-commands"
    identifier: eosc_tx
    weight: 40
weight: 40
draft: false
aliases: []
toc: false
auto_content: true
---

```
Usage:
  eosc tx [command]

Available Commands:
  create      Create a transaction with a single action
  push        Push a signed transaction to the chain.  Must be done online.
  sign        Sign a transaction produced by --write-transaction and submit it to the chain (unless --write-transaction is passed again).
  unpack      Unpack a transaction produced by --write-transaction and display all its actions (for review).  This does not submit anything to the chain.

Flags:
  -h, --help   help for tx
```
