// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/rex"
	"github.com/spf13/cobra"
)

var rexToSavings = &cobra.Command{
	Use:   "to-savings [account] [quantity]",
	Short: "Deposit REX tokens into your savings bucket.",
	Long:  "Deposit REX tokens into your savings bucket from your REX fund.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		account := toAccount(args[0], "account")
		quantity := toREXAsset(args[1], "quantity")

		pushEOSCActions(getAPI(), rex.NewMoveToSavings(
			account,
			quantity,
		))
	},
}

func init() {
	rexCmd.AddCommand(rexToSavings)
}
