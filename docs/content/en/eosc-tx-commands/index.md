---
title: "`eosc tx` Commands"
linktitle: "`eosc tx`"
description:
date: 2018-09-05
publishdate: 2018-09-05
lastmod: 2018-09-05
keywords: []
menu:
  docs:
    parent: "eosc-tx-commands"
    weight: 30
weight: 30
sections_weight: 30
draft: false
aliases: []
toc: true


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
