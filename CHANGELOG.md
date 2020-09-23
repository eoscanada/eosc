# Change log

The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this
project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html). See
[MAINTAINERS.md](./MAINTAINERS.md) for instructions to keep up to
date.

# Unreleased

### Added

* Added the `vault passwd` command, to change the passphrase on a
  passphrase-secured vault.

* A _short-form authority_ expression to specify authorities in `system updateauth` and `system newaccount` commands. Syntax:

```
- An optional threshold for the whole structure: "3=" (defauts to "1=")
- Comma-separated permission levels:
    - a public key
      or:
    - an account name, with optional "@permission" (defaults to "@active")
  For each permission levels, an optional "+2" suffix (defaults to "+1")

EXAMPLES

An authority with a threshold of 1, gated by a single key with a
weight of 1:

    EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV

An authority with a threshold of 1, gated by two accounts each with a weight of 1:

    myaccount,youraccount

An authority with a threshold of 2, gated by two accounts each with a weight of 1

    2=myaccount,youraccount

An authority with a threshold of 3, requiring admin (+2) and one of the two
employees (each +1):

    3=admin+2,employee1,employee2

An authority with a threshold of 3, gated by a key with a weight of 2, an
account with a weight of 3, and another account with a weight of 1:

    3=EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV+2,myaccount@secureperm+3,youraccount
```

this syntax is used to specify `system updateauth` and `system
newaccount` authority structures.  It is also available in the `cli`
package of the `eosc` library.

### Changed

* The `--auth-key` and `--auth-file` flags of the `system newaccount` command were removed, and changed to positional arguemnts, using the short-form autority syntax. The auth file support is removed, as the short-form syntax absorbs all of its functionality, in a much more terse way.

* The `--setpriv` flag on `system newaccount` was replaced by the `--additional-actions` option `setpriv`.

* `--sudo-wrap` now properly uses `eosio.wrap` as authorizer (was `eosio` and didn't work properly).
