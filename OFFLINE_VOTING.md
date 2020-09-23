How to vote using eosc and an airtight computer
-----------------------------------------------

[点击查看中文](./OFFLINE_VOTING-cn.md)

This step-by-step guide assumes you have a known account on mainnet
(referred to as `YOURACCOUNT`), that is provisioned with the public
key corresponding to the private key you'll be using in this process.

1. [ONLINE COMPUTER] Download a version of `eosc` from
   https://github.com/eoscanada/eosc/releases targeting your Airtight
   Computer's platform: Windows, Mac or Linux (a small footprint Linux
   distribution is preferred).

2. [ONLINE COMPUTER] Download and extract `eosc` for the Online
   Computer too.

3. [AIRTIGHT COMPUTER] Move the `eosc` distribution to your Airtight
   Computer and extract it over there.

4. [AIRTIGHT COMPUTER] Create a new Vault and punch in your private
   key(s):

```
eosc vault create --import --comment "Airtight wallet"
```

This will create a new `eosc-vault.json` file holding your key(s),
encrypted with your passphrase, on the Airtight Computer's disk.

5. [ONLINE COMPUTER] Run:

```
eosc get info
```

and copy its contents to the Airtight computer. What is really needed
is the `chain_id` and the `head_block_id` fields.

6. [AIRTIGHT COMPUTER] Create and sign your transaction:

```
eosc vote producers YOURACCOUNT PRODUCER1 PRODUCER2 PRODUCER3... --write-transaction signed-tx.json --expiration 3600 --offline-chain-id [CHAIN_ID] --offline-head-block [HEAD_BLOCK_ID] --offline-sign-key EOS6tg.....
```

This will create a file `signed-tx.json` containing the `voteproducer`
action, it will sign it with the local vault (prompt for the
passphrase) and set an expiration for in an hour.

You can check its unpacked data with `eosc tx unpack unsigned-tx.json`.

NOTE: you can use any `eosc` commands that produce transactions (like
`eosc transfer`, `eosc system delegatebw`) in the same way.


6. [ONLINE COMPUTER] Copy the `signed-tx.json` file to your Airtight
   Computer by whatever means.

7. [AIRTIGHT COMPUTER] Sign the transaction by replacing `EOS6tg...` by the Public Key that matches the private key in your Vault (and the YOURACCOUNT account authority):

```
eosc tx sign unsigned-tx.json --offline-sign-key EOS6tg..... --write-transaction signed-tx.json
```

This command will write `signed-tx.json` on disk.

8. [AIRTIGHT COMPUTER] Move the `signed-tx.json` back to the Online Computer.

9. [ONLINE COMPUTER] Push the signed transaction to the chain:

```
eosc tx push signed-tx.json
```

That's it!

### Troubleshooting

If the transaction failed for some reasons, you can read the error
message and see if they match one of:

* `UnAuthorized` errors: this means you have not signed the
  transaction with the keys required to authorize the
  `actor@permission` specified in the transaction's actions.

* `expired transaction`: you need to do the whole loop within the
  timeframe of the original `--expiration`. If you give `--expiration`
  more than an hour, note that you can only submit the transaction to
  the chain in the last hour of expiration.

* `Transaction's reference block did not match.`: this means maybe you
  didn't create the transaction from the same chain you're trying to
  push it to.

* `transaction bears irrelevant signatures from these keys: [\"EOS5J5....\"]`: this means you didn't need to sign with this key to authorize the `actor@permission` you specified in the transaction's actions.

* some other assertion errors, which are normal errors that would have
  occurred if you tried to sign it online, investigate with the
  contract in question.
