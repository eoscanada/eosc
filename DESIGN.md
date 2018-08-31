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
eosc multisig propose

eosc token transfer
