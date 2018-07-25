// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/spf13/cobra"
)

// vaultListCmd represents the list command
var vaultListCmd = &cobra.Command{
	Use:   "list",
	Short: "List public keys inside an eosc vault.",
	Long: `List public keys inside an eosc vault.

The wallet file contains a lits of public keys for easy reference, but
you cannot trust that these public keys have their counterpart in the
wallet, unless you check with the "list" command.
`,
	Run: func(cmd *cobra.Command, args []string) {
		vault := mustGetWallet()

		vault.PrintPublicKeys()
	},
}

func init() {
	vaultCmd.AddCommand(vaultListCmd)
}
