// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/rex"
	"github.com/spf13/cobra"
)

var rexClose = &cobra.Command{
	Use:   "close [account]",
	Short: "Removes all REX related entries from table.",
	Long:  "Free RAM from an account by removing its entry in the REX table. This action will fail if the account has any pending loans, refunds, or REX tokens.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		account := toAccount(args[0], "account")

		pushEOSCActions(getAPI(), rex.NewCloseREX(
			account,
		))
	},
}

func init() {
	rexCmd.AddCommand(rexClose)
}
