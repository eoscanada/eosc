---
title: "`eosc vote` Commands"
linktitle: "`eosc vote`"
description: Vote for Block Producers with eosc
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
Enter passphrase to unlock vault:
Voter [youraccount] voting for: [eoscanadacom]
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
