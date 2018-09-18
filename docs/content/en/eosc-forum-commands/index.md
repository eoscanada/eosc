---
title: "`eosc forum` Commands"
linktitle: "`eosc forum`"
description:
date: 2018-09-05
publishdate: 2018-09-05
lastmod: 2018-09-05
keywords: []
menu:
  docs:
    parent: "eosc-forum-commands"
    weight: 30
weight: 30
sections_weight: 30
draft: false
aliases: []
toc: true


---

### Forum

```
Usage:
  eosc forum [command]

Available Commands:
  post        Post a message
  propose     Submit a proposition for votes
  tally-votes Tally votes according to the `type` of the proposal.
  vote        Submit a vote from [voter] on [proposer]'s [proposal_name] with a [vote_value] agreed in the proposition.

Flags:
  -h, --help                     help for forum
      --target-contract string   Target account hosting the eosio.forum code (default "eosforumdapp")
```

