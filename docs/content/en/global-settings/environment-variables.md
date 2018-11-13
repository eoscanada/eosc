---
title: Environment Variable
linktitle: Environment Variable
description:
date: 2018-09-27
publishdate: 2018-09-27
lastmod: 2018-09-27
categories: [global setting]
keywords: []
authors: ["EOS Canada"]
menu:
  docs:
    parent: "global-settings"
    identifier: global-setting-environment-variables
    weight: 40
weight: 40
sections_weight: 40
draft: false
aliases: []
toc: false
---

Operating system environment variables can now automatically set any command-line flag:

* All global flags (those you get from eosc --help) can be set with the following pattern: EOSC_GLOBAL_FLAG_NAME.
   * For instance, the --wallet-url maps to EOSC_GLOBAL_WALLET_URL, etc.
* Any command and subcommands flags map to the following pattern: EOSC_COMMAND_SUBCOMMAND_CMD_FLAG_NAME.
   * For instance, eosc forum post --reply-to Mommy would be set as EOSC_FORUM_POST_CMD_REPLY_TO=Mommy. The CMD token separates command name from its flags. 
