---
title: eosc 0.7.9 更新
linktitle: eosc 0.7.9 更新
description: eosc v0.7.9 功能更新
date: 2018-09-27
publishdate: 2018-09-06
lastmod: 2018-09-27
keywords: []
weight: 40
sections_weight: 40
draft: false
aliases: []
toc: false
categories: [blog]
---

If you have been testing out different wallets to interact with the EOS blockchain, make sure to also check out `eosc`. `eosc` is always getting new features added by the team at EOS Canada. Here's a quick look into what has recently been released:

**v0.7.9 Features**

* `eosc get account` has been reformatted for easier reading. Displaying more info than `cleos`. For the old JSON blob, use the flag --json
* eosc tx `create`, `push`, `sign`, and `unpack` have been added for increased functionality for what you can do with a transaction
* `eosc multisig review` has been improved, plus added `setalimits` analysis
* The chain freeze tool has been improved, making it more precise
* `eosc system updateauth` added for permissions updating
* Offline signing flow has been improved for easier use - read our [walkthrough](https://github.com/eoscanada/eosc/blob/master/OFFLINE_VOTING.md)
