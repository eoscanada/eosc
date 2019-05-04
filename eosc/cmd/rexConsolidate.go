// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/rex"
	"github.com/spf13/cobra"
)

var rexConsolidate = &cobra.Command{
	Use:   "consolidate [account]",
	Short: "Consolidates any active REX maturity buckets.",
	Long:  "Consolidates any active REX maturity buckets into a single bucket that will mature in 4 days.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		account := toAccount(args[0], "account")

		pushEOSCActions(getAPI(), rex.NewConsolidate(
			account,
		))
	},
}

func init() {
	rexCmd.AddCommand(rexConsolidate)
}
