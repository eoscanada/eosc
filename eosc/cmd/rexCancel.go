// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/rex"
	"github.com/spf13/cobra"
)

var rexCancel = &cobra.Command{
	Use:   "cancel [account]",
	Short: "Cancels any unfilled sell orders.",
	Long:  "Cancels any unfilled sell orders for REX tokens.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		account := toAccount(args[0], "account")

		pushEOSCActions(getAPI(), rex.NewCancelREXOrder(
			account,
		))
	},
}

func init() {
	rexCmd.AddCommand(rexCancel)
}
