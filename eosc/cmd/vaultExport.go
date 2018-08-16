// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/spf13/cobra"
)

// vaultExportCommand represents the export command
var vaultExportCommand = &cobra.Command{
	Use:   "export",
	Short: "Export private keys (and corresponding public keys) inside an eosc vault.",
	Run: func(cmd *cobra.Command, args []string) {
		vault := mustGetWallet()

		vault.PrintPrivateKeys()
	},
}

func init() {
	vaultCmd.AddCommand(vaultExportCommand)
}
