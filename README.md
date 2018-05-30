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

The `eosc vault` is a simple yet powerful EOS wallet.


Sample usage
------------

    $ eosc vault create --keys 1 --shares 2/3
    EOS Vault starting
    "eosc-vault.json" already exists, manage that before running `eosc vault create`.

    $ eosc vault create --keys 1 --shares 2/3
    Creating 1 new keypair(s):
    - EOS123123
    Encrypting with 3 shares, 2 required for unlocking.
    Enter password for share 1:
    Enter password for share 2:
    Enter password for share 3:
    Encrypting to "eosc-vault.json"  done.  Please back that up.

    $ eosc vault serve
    Opening "eosc-vault.json".
    Public keys stored within:
    - EOS123123

    Unlock: 2/3 shares required (order doesn't matter)
    Enter first password:
    Enter second password:
    Listening on port :6666

    $ eosc vault serve -t 30s
    ...
    Listening on port :6666

    Locked after timeout, unlock with: 7823
    Enter unlock code: 7823
    Listening on port :6666

    $ eosc vault list
    Opening "eosc-vault.json".
    Public keys stored within:
    - EOS123123
    - EOS234234
    - EOS345345

Planned:

    $ eosc vault merge wallet1.json wallet2.json --shares 2/3
    $ eosc vault split EOS123123123 --shares 2/3
    $ eosc vault export
    Opening "eosc-vault.json".
    Public keys stored within:
    - EOS123123
    - EOS234234
    - EOS345345
    * 2/3 shares required to unlock (order doesn't matter)
    Enter first password:
    Enter second password:


Other features:

Support `--accept` which will ask you to confirm any signature on the command-line.
*  It will print what it can print..



FAQ
---

Why not use `cleos` instead ?

* hard to compile
* uses sha512 as key derivation, Argon2 is stronger.
* not work on Windows
* no need to run `keosd` separately and manage your keys on a Linux or Mac machine.
* single binary that brings `keosd` and `cleos` together
