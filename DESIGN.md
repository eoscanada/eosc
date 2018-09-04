## Basic help

Basic blockchain interaction and query:
  get      Read information from the blockchain, like accounts and balances.
  vote     Cast your votes, securely and simply.

Main contracts interfaces:
  system   System operations, like creating an account, setting up permission, etc.
  token    Transfer tokens from account to account

Built-in wallet management:
  vault    Built-in wallet, for in-memory signing, and keosd drop-in replacement.

Other tools:
  tools    Chain-freeze and other more involved tools


General tool layout
-------------------

eosc vault create
eosc vault create --import
eosc vault serve
eosc vault add

eosc vote producers
eosc vote list
eosc vote proxy

eosc get account
eosc get transaction
eosc get transactions
eosc get servants
eosc get code
eosc get abi
eosc get table
eosc get info
eosc get block
eosc get actions
eosc get balance [account]  - better than `get currency balance eosio.token [account]`  darn!

#eosc set code
#eosc set abi
#eosc set account
#eosc set account permission
#eosc set action permission

#eosc push action [contract] [action] payload -p permission@active
#eosc push transaction
#eosc push transactions

eosc tx create [[contract] [action] payload] -p permission@active
eosc tx create [path/to/actions.yaml] -p permission@active
eosc tx sign
eosc tx push (one, many trx, etc..)
eosc tx pack
eosc tx unpack
eosc tx pack-actions
eosc tx unpack-actions
eosc tx review
eosc tx combine [write_to.json] [tx1.json] [tx2.json] [tx3.json]
eosc tx combine [write_to.json] [tx1.json] [tx2.json] [tx3.json]
eosc tx combine [write_to.json] [tx1.json] [tx2.json] [tx3.json]

eosc system delegatebw
eosc system newaccount
eosc system sellram
eosc system buyram
eosc system setcode
eosc system setabi
eosc system updateauth
eosc system linkauth
eosc system voteproducers

eosc multisig propose
eosc multisig list proposer
eosc multisig approvals proposer proposal
eosc multisig review // LIST those requested for approval (and their approval state?)

eosc token transfer


----------------------

eosc compile
eosc boot
eosc test



eosc build skeleton --advanced forum
------------------------------------

* write a helloworld contract, simple use of a database
* build make / build.sh
* sample ABI with a struct
* .eosc-test.yaml
* .eosc-build.yaml


eosc build skeleton --simple forum
----------------------------------

* simple helloworld cpp, hpp
* sample test suite under `tests` ?


eosc build env
--------------

Drops you in the build environment in `docker`

Detect `docker` install, and gives you a pointer to how to install it
or something.


eosc build [now|all]
------------

* `eosc build` or `eosc build now` are equal
* `eosc build all` navigates all subdirectories for `.eosc-build.yaml`

.eosc-build.yaml
---
v: 1
wasmsdk: 1.2.3
build: cmake
---

eosc test [--keep]
---------

.eosc-test.yaml
---
nodeos_version: 1.2.2
nodeos_config:
  verbose-http-error: 1.2.3.4
boot_sequence: mainnet-1.2.2
vault_keys:
- 5K3434242342342412312313123123  ; Public = EOS234234
- 5K123123123213123121212121212121  ; Public = EOS123123
data_dir: ./test/data
script: ./run.sh
---

* eosc nodeos run, pick up some ports, pass in some config (from `eosc-test.yaml`
* eosc test
  * runs a vault with `New Vault`, and boots the chain with that PrivKey, belongs to EOSIO
  * vault litsen on :6668, sets
* run `run.sh`, set ENV: `EOSC_GLOBAL_WALLET_URL=http://localhost:6668


run.sh
---
#!/bin/bash -e

# Has EOSC_GLOBAL_API_URL
# Has EOSC_GLOBAL_WALLET_URL
# Has EOSC_GLOBAL_MANAGEOS_URL

if [ $IN_CIRCLE_CI ] {
  eosc pitreos restore ./nodeos-data -t lastbomb
  eosc nodeos start
} else {
  eosc nodeos start
  eosc nodeos bios-boot mainnet.yaml
}

eosc tx sign ./fixtures/megatx_test.json

eosc system newaccount
eosc system setcontract bob.wasm bob.abi

eosc tx create eosio setprods '{"prods": "eoscanadacom"}'
eosc get table eosio eosio producers -L 123 | jq -r .1.name

eosc nodeos production pause
eosc nodeos production resume
eosc nodeos stop

BLOCKNUM=$(eosc get info | jq .block_num)
eosc pitreos backup ./nodeos-data -t lastbomb -meta '{"block_id": $BLOCKNUM}'

eosc nodeos run (through docker)

eosc assert table-contents eosio producers ./table_contents.json


# Implicit teardown.  Shutdown `nodeos`, you can restart it with `eosc nodeos run`
---
