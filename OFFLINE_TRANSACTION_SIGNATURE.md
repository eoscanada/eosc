## Offline transaction signature


`eosc` has everything needed for managing an airtight offline wallet.
Here is a sample use. From within `/tmp/airtight`, run:

* `eosc vault create --import --comment "Airtight wallet"`

which will output:
```
PLEASE READ:
We are now going to ask you to paste your private keys, one at a time.
They will not be shown on screen.
Please verify that the public keys printed on screen correspond to what you have noted

Paste your first private key:       <em>[PASTE PRIVATE KEY HERE]</em>
- Scanned private key corresponding to EOS5J53bNtH5H1bJTyP237fr5LF6eQSQohfGY5iMCgpC4HpXApJBr
Paste your next private key or hit ENTER if you are done:           <em>[PASTE SECOND PRIVATE KEY HERE]</em>
- Scanned private key corresponding to EOS6tgsdv6S7N1GYWgX8QEBAsanAwXwuaEkv11GGtteyk5ELqSzVP
Paste your next private key or hit ENTER if you are done:           <em>[HIT ENTER HERE]</em>
Imported 2 keys.

You will be asked to provide a passphrase to secure your newly created vault.
Make sure you make it long and strong.

Enter passphrase to encrypt your vault:       <em>[ENTER YOUR PASSPHRASE HERE]</em>
Confirm passphrase:                           <em>[RE-ENTER YOUR PASSPHRASE HERE]</em>

Wallet file "./eosc-vault.json" written to disk.
Here are the keys that were ADDED during this operation (use `list` to see them all):
- EOS5J53bNtH5H1bJTyP237fr5LF6eQSQohfGY5iMCgpC4HpXApJBr
- EOS6tgsdv6S7N1GYWgX8QEBAsanAwXwuaEkv11GGtteyk5ELqSzVP
Total keys stored: 2
```

This operation creates your `eosc-vault.json` wallet.  It can be done
completely offline as it does not touch the network at all.

Now, let's try using it. Issue those commands from a computer with internet access (which in turn does not have access to your eosc vault wallet):

* `eosc transfer airtight11111 cancancan123 1.0000 -m "Can't say I haven't paid you now" --write-transaction transaction.json --skip-sign --expiration 3600`

This will output:

```
{
  "expiration": "2018-08-27T04:18:43",
  "ref_block_num": 43579,
  "ref_block_prefix": 3787183056,
  "max_net_usage_words": 0,
  "max_cpu_usage_ms": 0,
  "delay_sec": 0,
  "context_free_actions": [],
  "actions": [
    {
      "account": "eosio.token",
      "name": "transfer",
      "authorization": [
        {
          "actor": "airtight1111",
          "permission": "active"
        }
      ],
      "data": {
        "from": "airtight1111",
        "to": "cancancan123",
        "quantity": "0.0001 EOS",
        "memo": "Can't say I haven't paid you now"
      }
    }
  ],
  "transaction_extensions": [],
  "signatures": [],
  "context_free_data": []
}
---
Transaction written to "../transaction.json"
Sign offline with: --offline-chain-id=5fff1dae8dc8e2fc4d5b23b2c7665c97f9e9d8edf2b6485a86ba311c25639191
Above is a pretty-printed representation of the outputted file
```

Now copy the `transaction.json` file to another computer. Note the `--offline-chain-id=...` output above.

On the airtight computer, run:

* `eosc tx sign path/to/transaction.json --offline-sign-key EOS6tgsdv6S7N1GYWgX8QEBAsanAwXwuaEkv11GGtteyk5ELqSzVP  --offline-chain-id=5fff1dae8dc8e2fc4d5b23b2c7665c97f9e9d8edf2b6485a86ba311c25639191 --write-transaction signedtx.json`

You can enter as many `--offline-sign-key` as you want, or separate each entry with a `,` (comma), like this: `--offline-sign-key EOS6tgsdv6S7N1GYWgX8QEBAsanAwXwuaEkv11GGtteyk5ELqSzVP,EOS5J53bNtH5H1bJTyP237fr5LF6eQSQohfGY5iMCgpC4HpXApJBr`.  This will output:

```
Enter passphrase to decrypt your vault:
{
  "expiration": "2018-08-27T04:18:43",
  "ref_block_num": 43579,
  "ref_block_prefix": 3787183056,
  "max_net_usage_words": 0,
  "max_cpu_usage_ms": 0,
  "delay_sec": 0,
  "context_free_actions": [],
  "actions": [
    {
      "account": "eosio.token",
      "name": "transfer",
      "authorization": [
        {
          "actor": "airtight1111",
          "permission": "active"
        }
      ],
      "data": "104208b93197af33304498064d83a641010000000000000004454f53000000002043616e277420736179204920686176656e2774207061696420796f75206e6f77"
    }
  ],
  "transaction_extensions": [],
  "signatures": [
    "SIG_K1_KbsvdjfXTzbn7DaghCbFyaxXnv9BHJaAVWZkeJ3AbvnznNgipSLj7BbZajZ1ECZEWKSMFtMbuqfUWQGq2tDzW2n7Fz5KTV",
    "SIG_K1_K24N3kUfTLMJGwUUWV2qNiyhxAJsbVtpXh2KUj3SHGLTivPBru46QYX9v9gQj7G2yKHp6eZ786hHfJwuzjeZFF7atPfpTY"
  ],
  "context_free_data": []
}
---
Transaction written to "signedtx.json"
Sign offline with: --offline-chain-id=5fff1dae8dc8e2fc4d5b23b2c7665c97f9e9d8edf2b6485a86ba311c25639191
Above is a pretty-printed representation of the outputted file
```

You are now ready to bring `signedtx.json` back to the online computer, and run:

* `eosc tx push signedtx.json`

This will succeed, in which case you're done, or fail with some of these:

* `transaction bears irrelevant signatures from these keys: [\"EOS5J53bNtH5H1bJTyP237fr5LF6eQSQohfGY5iMCgpC4HpXApJBr\"]`: this means you didn't need to sign with this key to authorize the `actor@permission` you specified in the transaction's actions.

* `UnAuthorized` errors: this means you have not signed the
  transaction with the keys required to authorize the
  `actor@permission` specified in the transaction's actions.

* `Transaction's reference block did not match.`: this means you
  didn't create the transaction from the same chain you're trying to
  push it to.

* `expired transaction`: you need to do the whole loop within the
  timeframe of the original `--expiration`. If you give `--expiration`
  more than an hour, note that you can only submit the transaction to
  the chain in the last hour of expiration.

* some other assertion errors, which are normal errors that would have
  occurred if you tried to sign it online, investigate with the
  contract in question.



### Offline multisig

For added security with offline transaction signing, consider using
the `multisig` facility of the EOS blockchain.

It is possible to propose a transaction on-chain (through `eosc
multisig propose`), reviewable by all, and prepare an offline
transaction simply to approve that transaction.  Using this method,
you could have X number of airtight computers that each provide one
piece of the puzzle, controlled by different individuals, without the
need to stitch transaction files, or gather signatures into a single
place.
