How to vote using eosc and an airtight computer
-----------------------------------------------

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

5. [ONLINE COMPUTER] Prepare your `voteproducer` transaction with:

```
eosc vote producers YOURACCOUNT PRODUCER1 PRODUCER2 PRODUCER3... --write-transaction unsigned-tx.json --skip-sign --expiration 3600
```

This will create a file `unsigned-tx.json` containing the
`voteproducer` action. It will not try to sign it (that will happen on
the Airtight computer) and sets an expiration of one hour for the
roundtrip to happen.

You can check its unpacked data with `eosc tx unpack unsigned-tx.json`.

NOTE: you can use any `eosc` commands that produce transactions (like
`eosc transfer`, `eosc system delegatebw`) in the same way.


6. [ONLINE COMPUTER] Copy the `unsigned-tx.json` file to your Airtight
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
