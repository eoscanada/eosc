# Use `eosc` to Set a Proxy Account

## Use `eosc` to Create a New Account

Grab the latest version of `eosc` from https://github.com/eoscanada/eosc/releases

Once unpacked, run:
```
eosc vault create --keys [int]  // if you want eosc to create keys for you

OR

eosc vault create --import [PUBLIC KEY]  // if you have a public key/private key pair that you would like to encrypt into the vault
```
By default, the vault file will be "./eosc-vault.json"
To assign a different name, use the flag `--vault-file [VAULT NAME]`

You will prompted to create a passphrase, and then to confirm that passphrase.
Make sure to keep a backup of your vault file, especially if you use eosc to generate the private key for you.

Now that we have a public key/private key pair, we'll need to assign it to a new account.
The design of EOS requires you to have an account to create an account. So we'll need 
to use an already existing account on the EOS mainnet to create this new account.

If your already existing account does not have an eosc vault file already,
you can create one using the above 
`eosc vault create --import`.

To create the new account, run the command:
```
eosc system newaccount [EXISTING ACCOUNT NAME] [NEW ACCOUNT NAME]
```
You'll want to attach all necessary flags to this. To note are:
```
Flags:
      --auth-file string     File containing owner and active permissions authorities. See example in --help
      --auth-key string      Public key to use for both owner and active permissions.
      --buy-ram string       The amount of EOS to spend to buy RAM for the new account (at current EOS/RAM market price)
      --buy-ram-kbytes int   The amount of RAM kibibytes (KiB) to purchase for the new account.  Defaults to 8 KiB. (default 8)
  -h, --help                 help for newaccount
      --setpriv              Make this account a privileged account (reserved to the 'eosio' system account)
      --stake-cpu string     Amount of EOS to stake for CPU bandwidth (required)
      --stake-net string     Amount of EOS to stake for Network bandwidth (required)
      --transfer             Transfer voting power and right to unstake EOS to receiver
```
The new account will require RAM. The current minimum is 3 KiB. By default eosc will purchase 8 KiB unless using 
either the `--buy-ram` or `--buy-ram-kbytes` flags.
You will now have the new account registered on the mainnet.

## How to Register Your New Account as a Proxy

To register [NEW ACCOUNT] as a proxy, run this command:
```
eosc system regproxy [ACCOUNT NAME]
```
Your account will now be registered as a proxy. 
If you would like to delegate the voting strength from [EXISTING ACCOUNT]
to [NEW ACCOUNT] (now a proxy), you would run the command:
```
eosc vote proxy [VOTING ACCOUNT] [NEW ACCOUNT]
```
If you want to verify on-chain if your proxy registration has gone through, you can run the command
```
eosc get account [ACCOUNT NAME] --json
```
You will look for the field `is_proxy:`

If the value is 0, then the account is not registered as a proxy. If it is 1, then the registration was successful.

## Voting With Your Proxy Account

To vote with your proxy, you follow the same steps that any account would for voting.
```
eosc vote producers [ACCOUNT NAME] [Producer1] [Producer2] ...
```
An account can vote for up to 30 Block Producers at a time.
Some handy tools that we've built into `eosc` is the ability to call
`eosc vote status [ACCOUNT NAME]`
in case you wanted to copy the voting string of another account without
having to type in all of the producer names manually.
For recasting the same vote to keep your [vote strength up](https://www.eoscanada.com/en/how-is-your-vote-strength-calculated-on-eos), you can use
`eosc vote recast [ACCOUNT NAME]`

## Voting For a Proxy Account

To vote for an account that has been registered as a proxy, run the command:
```
eosc vote proxy [YOUR ACCOUNT] [PROXY ACCOUNT]
```
To verify that your vote has been cast for that proxy, you can run the command
```
eosc get account [YOUR ACCOUNT] --json
```
and you will be see at the end of the json blob
```
"voter_info": {
    "owner": "YOUR ACCOUNT",
    "proxy": "PROXY ACCOUNT",
    "producers": [],
    "staked": ,
    "last_vote_weight": ,
    "proxied_vote_weight": 0,
    "is_proxy": 0
```
