`eosc` command-line swiss-knife
-------------------------------

`eosc` is a cross-platform (Windows, Mac and Linux) command-line tool
for interacting with an EOS.IO blockchain.

It contains tools for voting and a Vault to securely store private
keys.

It is based on the `eos-go` library, used in `eos-bios`, the tool used
to launch the EOS mainnet.


eosc vault
----------

The `eosc` Vault is a simple yet powerful EOS wallet.


Sample usage
------------

    $ eos-vault create --keys 1 --shares 2/3
    EOS Vault starting
    "eos-vault-wallet.json" already exists, manage that before running `eos-vault create`.

    $ eos-vault create --keys 1 --shares 2/3
    Creating 1 new keypair(s):
    - EOS123123
    Encrypting with 3 shares, 2 required for unlocking.
    Enter password for share 1:
    Enter password for share 2:
    Enter password for share 3:
    Encrypting to "eos-vault-wallet.json"  done.  Please back that up.

    $ eos-vault serve
    Opening "eos-vault-wallet.json".
    Public keys stored within:
    - EOS123123

    Unlock: 2/3 shares required (order doesn't matter)
    Enter first password:
    Enter second password:
    Listening on port :6666

    $ eos-vault serve -t 30s
    ...
    Listening on port :6666

    Locked after timeout, unlock with: 7823
    Enter unlock code: 7823
    Listening on port :6666

    $ eos-vault list
    Opening "eos-vault-wallet.json".
    Public keys stored within:
    - EOS123123
    - EOS234234
    - EOS345345

    $ eos-vault export
    Opening "eos-vault-wallet.json".
    Public keys stored within:
    - EOS123123
    - EOS234234
    - EOS345345
    * 2/3 shares required to unlock (order doesn't matter)
    Enter first password:
    Enter second password:

Planned:

    $ eos-vault merge wallet1.json wallet2.json --shares 2/3
    $ eos-vault split EOS123123123 --shares 2/3

Other features
--------------

Support `--accept` which will ask you to confirm any signature on the command-line.
  It will print what it can print..
