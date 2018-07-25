// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/spf13/cobra"
)

// vaultCmd represents the vault command
var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "The eosc vault is a secure EOS key vault and a wallet server",
	Long:  `It is a drop-in replacement for keosd with additional features.`,
}

func init() {
	RootCmd.AddCommand(vaultCmd)
}
