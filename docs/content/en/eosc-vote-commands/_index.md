---
title: vote Commands Overview
linktitle: vote Commands Overview
description: Overview of the `vote` section of eosc
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [vote]
keywords: [usage,docs]
menu:
  docs:
    parent: "eosc-vote-commands"
    identifier: eosc_vote
    weight: 1
weight: 0001	#rem
draft: false
aliases: [/overview/introduction/]
toc: false
auto_content: true
---

```
Usage:
  eosc vote [command]

Available Commands:
  cancel-all     Cancel all votes currently casted for producers and reset proxy voting.
  list-producers Retrieve the list of registered producers
  producers      Cast your vote for 1 to 30 producers. View them with 'list'
  proxy          Cast your vote for a proxy voter
  recast         Recast your vote for the same producers
  status         Display the current vote status for a given account.

Flags:
  -h, --help   help for vote
```

# Vote for Block Producers with eosc


Start by running:


```
eosc vote --help
```


You can then see all the instructions of the commands.



To start voting by running:


```
eosc vote producers [your account] [producer1] [producer2] [producer3]
```


You can vote for up to 30 Block Producers with each account, browse 
all block producer accounts on https://eosquery.com.



An interactive prompt would ask for your passphrase to proceed:


```
Enter passphrase to unlock vault: ***
Done
```


**If you haven't imported your privated keys into eosc vault, you 
wouldn't have the passphrase, see how to set it up [here]()**

This will sign your vote transaction locally, and submit the
transaction to the network through the `https://mainnet.eoscanada.com`
endpoint. You can also point to some other endpoints that are on the
main network with the flag `-u` or `--api-url`.


Make sure you have eosc version `v0.7.0` or higher, you can run the 
following command to see the version you installed:


```
eosc version
```


## Vote for Proxy Accounts


run:


```
eosc vote proxy [your account] [proxy name]
```
