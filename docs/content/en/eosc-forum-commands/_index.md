---
title: forum Commands Overview
linktitle: forum Commands Overview
description: Overview of the `forum` section of eosc
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [forum]
keywords: [usage,docs]
menu:
  docs:
    parent: "eosc-forum-commands"
    identifier: eosc_forum
    weight: 1
weight: 0001	#rem
draft: false
aliases: [/overview/introduction/]
toc: false
auto_content: true
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

