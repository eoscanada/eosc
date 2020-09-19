// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"context"
	"fmt"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
	"github.com/eoscanada/eosc/cli"
	"github.com/spf13/cobra"
)

func init() {
	systemCmd.AddCommand(systemUpdateauthCmd)
}

var systemUpdateauthCmd = &cobra.Command{
	Use:   `updateauth [account] [permission_name] [parent permission or ""] [authority]`,
	Short: "Set or update a permission on an account. See --help for more details.",
	Long: `Set or update a permission on an account.

The [authority] field is expressed using this short-form syntax:
` + shortFormAuthHelp,
	Args: cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		account := toAccount(args[0], "account")
		permissionName := toName(args[1], "permission_name")

		var parent eos.Name
		if args[2] != "" {
			parent = toName(args[2], "parent permission")
		}
		authParam := args[3]

		auth, err := cli.ParseShortFormAuth(authParam)
		errorCheck("parsing authority", err)
		err = ValidateAuth(auth)
		errorCheck("validating authority", err)

		api := getAPI()

		var updateAuthActionPermission = "active"
		if parent == "" {
			updateAuthActionPermission = "owner"
		}

		pushEOSCActions(
			context.Background(),
			api,
			system.NewUpdateAuth(
				account,
				eos.PermissionName(permissionName),
				eos.PermissionName(parent),
				*auth,
				eos.PermissionName(updateAuthActionPermission),
			),
		)
	},
}

func ValidateAuth(auth *eos.Authority) error {
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
