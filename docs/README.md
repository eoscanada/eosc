# eosc docs

This website holds documentation for the different parts of `eosc`.

Contributions are welcome.  See CONTRIBUTING.md at the root of this
repository.

## Build

To view the documentation site locally, you need to clone this repository:

```bash
git clone https://github.com/eoscanada/eosc.git
```

Get the latest [Hugo](https://github.com/gohugoio/hugo) release and build the docs with:

```bash
▶ hugo server

Started building sites ...
.
.
Serving pages from memory
Web Server is available at http://localhost:1313/ (bind address 127.0.0.1)
Press Ctrl+C to stop
```

Then open your browser at http://localhost:1313/


## Draft Sections


About eosc
* What is eosc
* Features overview
* Benefits
* License

Getting Started
* Overview
* Quick Start
* Install eosc
* Basic Usage

Global settings
* Overview
* Global flags
* Environment variables

---

`eosc get` commands
* Overview
* account     retrieve account information for a given name
* balance     Retrieve currency balance for an account
* block       Get block data at a given height, or directly with a block hash
* code        retrieve account information for a given name
* info        Retrieve blockchain infos, like head block, chain ID, etc..
* table       Fetch data from a table in a contract on chain

`eosc vault` commands
*

`eosc system` commands
* Overview
* `bidname`      Bid on a premium account name
* `buyrambytes`  Buy RAM at market price, for a given number of bytes.
* claimrewards Claim block production rewards. Once per day, don't forget it!
* delegatebw   Delegate some CPU and Network bandwidth, to yourself or others.
* linkauth     Assign a permission to the given code::action pair
* newaccount   Create a new account
* regproducer  Register an account as a block producer candidate
* sellram      Sell the [num bytes] amount of bytes of RAM on the RAM market.
* setabi       Set ABI only on an account
* setcode      Set code only on an account
* setcontract  Set both code and ABI on an account
* undelegatebw Undelegate some CPU and Network bandwidth.
* unregprod    Unregister producer account temporarily.
* updateauth   Set or update a permission on an account. See --help for more details.

`eosc multisig` commands
* Overview
* approve     Approve a transaction in the eosio.msig contract
* cancel      Cancel a transaction in the eosio.msig contract
* exec        Execute a transaction in the eosio.msig contract
* propose     Propose a new transaction in the eosio.msig contract
* review      Review a proposal in the eosio.msig contract
* unapprove   Unapprove a transaction in the eosio.msig contract

`eosc forum` commands
* Overview
* post        Post a message
* propose     Submit a proposition for votes
* tally-votes Tally votes according to the `type` of the proposal.
* vote        Submit a vote from [voter] on [proposer]'s [proposal_name] with a [vote_value] agreed in the proposition.

`eosc tx` commands
* Overview
* create      Create a transaction with a single action
* push        Push a signed transaction to the chain.  Must be done online.
* sign        Sign a transaction produced by --write-transaction and submit it to the chain (unless --write-transaction is passed again).
* unpack      Unpack a transaction produced by --write-transaction and display all its actions (for review).  This does not submit anything to the chain.

`eosc vote` commands
* Overview
* list        Retrieve registered producers
* producers   Cast your vote for 1 to 30 producers. View them with 'list'
* proxy       Cast your vote for a proxy voter

`eosc tools` commands
* Overview
* chain-freeze    Runs a p2p protocol-level proxy, and stop sync'ing the chain at the given block-num.
* sell-account    Create a multisig transaction that both parties need to approve in order to do an atomic sale of your account.

---

Help
* FAQ         (FAQ-like content)
* Community   (Telegram)


-----------------------------

We generated help based on the `eosc` output, and produce
`generated_help.json` under `/data`.

---

.Data["eosc-tools-chain-freeze"]["short"]

```
{{ if $parent := .Data["eosc-tools-chain-freeze"]["parent"]; $parent }}
{{ with .Data[$parent]["persistent_flags"] }}Parent flags:
{{ . }}{{ end }}
{{ end }}
```

---

Shortcode `generated_help` that reads the `eosccommand` from the article, and produces
output based on the `data/generatedhelp.json` file.

---

Articles' manifest would include:

```
...
title: ...
...
eosccommand: eosc-tools
---

{% generated_help  %}
```
