---
title: Frequently Asked Questions
linktitle: FAQ
description:
date: 2018-09-27
categories: [help]
menu:
  docs:
    parent: "help"
    identifier: faq-help
    weight: 40
keywords: [faqs]
weight: 40
toc: false
aliases: []
---

### How do I vote?
Command: 
```bash
`eosc vote producers [YourAccountName] [producer1] [producer2] … [producer30]`
```
You can vote for up to 30 Block Producers. Make sure that you have the correct account name for the producer you wish to vote for. You can always run the command `eosc vote list` to view a complete list of all accounts that have called the `regproducer` action. You should always verify on a Block Producer’s website/social media that you are using the correct name, as some imposters have been found.

### How do I check my account info?
Command: 
```bash
`eosc get account [AccountName]`
```
What you’ll find: 
* Creation date
* RAM usage/quota
* CPU usage/quota
* Network usage/quota
* Account permissions
* total staked
* last vote weight
* Is this account a proxy? How much weight is being proxied?
* "deferred_trx_id": are there any deferred transactions associated to this account
* when the 72 hour countdown for unstaking was last started
* the amount currently being unstaked
* Note that you will not find the amount of unstaked EOS you hold in your account. To find this, you’ll need to run the command `eosc get balance [AccountName]`

### How do I transfer EOS to another account?
Command: 
```bash
`eosc transfer [SenderAccount] [ReceiverAccount] [amount]`
```
This command will send the amount of EOS specified in [amount] from the [SenderAccount] to the [ReceiverAccount]. You can put up to 4 decimal places at most in the [amount] field. There are 3 iterations to how you can write the [amount]: 1, 1.0000, or 1.0000 EOS.
If you’d like to attach a memo to the transaction, you can add `-m` or `--memo` to the command followed by the memo. 

### How do I delegate bandwidth?
Command: 
```bash
`eosc system delegatebw [from] [receiver] [network bw stake qty] [cpu bw stake qty]`
```
This command will delegate bandwidth from the [from] account to the [receiver] account. You can put whatever amount you’d like in the two quantity fields, they do not have to match. If you’d like to delegate only CPU and not Network, you would have to include a 0 for the Network field. 

A feature you can use if you’d like to transfer tokens to an account and have them be staked right away, is that you can add the flag `--transfer` to this action so that the tokens will begin staked for the other account, and the [receive] account will be able to unstake those tokens into their account.

### How do I undelegate bandwidth?
Command: 
```bash
`eosc system undelegatebw [from] [receiver] [network bw stake qty] [cpu bw stake qty]`
```
This command will undelegate bandwidth from the [from] account to the [receiver] account. You can put whatever amount you’d like in the two quantity fields, they do not have to match. If you’d like to undelegate only CPU and not Network, you would have to include a 0 for the Network field. 

Note that unstaking takes 72 hours before the tokens are fully available back to the account. 

### How do I buy RAM?
Command: 
```bash
`eosc system buyrambytes [payer] [receiver] [num bytes]`
```
You can specify an amount of bytes of RAM that you’d like to purchase for an account. 
Note that all RAM transactions will also charge a 0.5% fee to the amount. You can view the estimated cost for the amount you’d like to purchase at https://www.eosrp.io/ 
You will not be prompted with a price for purchase when using `eosc`, so please ensure that you have entered in the correct amount of bytes that you’d like.

### How do I sell RAM?
Command: 
```bash
`eosc system sellram [YourAccountName] [num bytes]`
```
You can specify the amount of bytes of RAM that you’d like to sell from your account. 
Note that all RAM transactions will also charge a 0.5% fee to the amount. 
You will not be prompted with a price for selling when using `eosc`, so please ensure that you have entered in the correct amount of bytes that you’d like to sell.

### How do I bid on a namespace?
Command: 
```bash
`eosc system bidname [bidder_account_name] [premium_account_name] [bid quantity]`
```
Bid on a premium account name/namespace. If you aren’t sure what this is for, please read our article which outlines the [rules for namespace auctions.](https://www.eoscanada.com/en/everything-you-need-to-know-about-namespace-bidding-about-eos) 

### How do I delegate my votes to a proxy?
Command: 
```bash
`eosc vote proxy [voter name] [proxy name]`
```
This will allow you to delegate your voting weight to a proxy account. If you are unsure if an account is listed as a proxy, you can check by using `eosc get account [account name]` and within the information listed, check to see if "is_proxy": is listed as `1` or `0`. `1` means that the account has registered as a proxy.

### How do I check current information of the chain?
Command: 
```bash
`eosc get info`
```
What you’ll find:
* chain ID
* Head block number
* Last Irreversible Block
* Current Producer
* "virtual_block_cpu_limit": 23744289,
* "virtual_block_net_limit": 1048576000,
* "block_cpu_limit": 200000,
* "block_net_limit": 1048576

### How Do I Create a New Account?
Command:
```bash
`eosc system newaccount [creator] [new_account_name] [flags]`
```
Creating a new account will require you to use the [flags] that are listed.
```bash
Flags:
      --auth-file string     File containing owner and active permissions authorities. See example in --help
      --auth-key string      Public key to use for both owner and active permissions.
      --buy-ram string       The amount of EOS to spend to buy RAM for the new account (at current EOS/RAM market price)
      --buy-ram-kbytes int   The amount of RAM kibibytes (KiB) to purchase for the new account.  Defaults to 8 KiB. (default 8)
  -h, --help                 help for newaccount
      --setpriv              Make this account a privileged account (reserved to the 'eosio' system account)
      --stake-cpu string     Amount of EOS to stake foor CPU bandwidth (required)
      --stake-net string     Amount of EOS to stake for Network bandwidth (required)
      --transfer             Transfer voting power and right to unstake EOS to receiver
``` 
`--auth-file string` is to be used if you have setup a vault already using `eosc`

`--auth-key string` is to be used if you would like to enter in the Public Key to be assigned for both the Owner and Active permissions.

`--buy-ram string` is to be used if you would like to spend a specific amount of EOS on buying RAM for an account

`--buy-ramkbytes int` is to be used when you want to purchase an amount of KiB of RAM for an account. If this is not used, the default of 8 KiB is used.

`--stake-cpu string` is required to be used. This is how you set a certain amount of EOS to be staked for CPU for the new account.

`--stake-net string` is required to be used. This is how you set a certain amount of EOS to be staked for Network for the new account.

`--transfer` is to be used if you would like to also allow the voting and unstaking rights of the tokens used for `--stake-cpu string` and `--stake-net string` to be transferred to the new account. 
