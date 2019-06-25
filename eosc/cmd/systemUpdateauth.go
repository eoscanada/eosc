// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"fmt"
	"sort"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var systemUpdateauthCmd = &cobra.Command{
	Use:   `updateauth [account] [permission_name] [parent permission or ""] [authority]`,
	Short: "Set or update a permission on an account. See --help for more details.",
	Long: `Set or update a permission on an account.

The [authority] field can be either a *public key* or a path to a YAML
file.

If you specify a public key, a simple 'authority' structure is built,
with a threshold of 1, and a single key.

Otherwise, it should be a path to a YAML file.  Here is a sample YAML
authority file:

---
threshold: 3
keys:
- key: EOS6MRyAjQq8ud7hVNYcfn................tHuGYqET5GDW5CV
  weight: 1
accounts:
- permission:
    actor: accountname
    permission: namedperm
  weight: 1
waits:
- wait_sec: 300
  weight: 1
---

`,
	Args: cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		account := toAccount(args[0], "account")
		permissionName := toName(args[1], "permission_name")

		var parent eos.Name
		if args[2] != "" {
			parent = toName(args[2], "parent permission")
		}
		authParam := args[3]

		var auth eos.Authority
		authKey, err := ecc.NewPublicKey(authParam)
		if err == nil {
			auth = eos.Authority{
				Threshold: 1,
				Keys: []eos.KeyWeight{
					{PublicKey: authKey, Weight: 1},
				},
			}
		} else {
			err := loadYAMLOrJSONFile(authParam, &auth)
			errorCheck("authority file invalid", err)
		}

		sortAuth(auth)

		err = ValidateAuth(auth)
		errorCheck("authority file invalid", err)

		api := getAPI()

		var updateAuthActionPermission = "active"
		if parent == "" {
			updateAuthActionPermission = "owner"
		}
		pushEOSCActions(api, system.NewUpdateAuth(account, eos.PermissionName(permissionName), eos.PermissionName(parent), auth, eos.PermissionName(updateAuthActionPermission)))
	},
}

func init() {
	systemCmd.AddCommand(systemUpdateauthCmd)
}

func ValidateAuth(auth eos.Authority) error {
	for idx, account := range auth.Accounts {
		if len(account.Permission.Permission) == 0 {
			return fmt.Errorf("account #%d missing permission", idx)
		}
		if len(account.Permission.Actor) == 0 {
			return fmt.Errorf("account #%d missing actor", idx)
		}

		if account.Weight == 0 {
			return fmt.Errorf("account #%d missing weight", idx)
		}
	}

	for idx, key := range auth.Keys {
		if len(key.PublicKey.Content) == 0 {
			return fmt.Errorf("key #%d missing publicKey", idx)
		}

		if key.Weight == 0 {
			return fmt.Errorf("key #%d missing weight", idx)

		}
	}

	for idx, wait := range auth.Waits {
		if wait.WaitSec == 0 {
			return fmt.Errorf("wait #%d cannot be 0", idx)
		}

		if wait.Weight == 0 {
			return fmt.Errorf("wait #%d is missing weight", idx)
		}
	}
	return nil
}

func sortAuth(auth eos.Authority) {
	sort.Slice(auth.Keys, func(i, j int) bool {
		return auth.Keys[i].PublicKey.String() < auth.Keys[j].PublicKey.String()
	})
	sort.Slice(auth.Accounts, func(i, j int) bool {
		perm1 := auth.Accounts[i].Permission
		perm2 := auth.Accounts[j].Permission
		if perm1.Actor == perm2.Actor {
			return perm1.Permission < perm2.Permission
		}
		return perm1.Actor < perm2.Actor
	})
	sort.Slice(auth.Waits, func(i, j int) bool {
		return auth.Waits[i].WaitSec < auth.Waits[j].WaitSec
	})
}
