// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/rex"
	"github.com/spf13/cobra"
)

var rexFromSavings = &cobra.Command{
	Use:   "from-savings [account] [quantity]",
	Short: "Withdraw REX tokens from your savings bucket.",
	Long:  "Withdraw REX tokens from your savings bucket into your REX fund. Those funds will become available in 4 days.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		account := toAccount(args[0], "account")
		quantity := toREXAsset(args[1], "quantity")

		pushEOSCActions(getAPI(), rex.NewMoveFromSavings(
			account,
			quantity,
		))
	},
}

func init() {
	rexCmd.AddCommand(rexFromSavings)
}
