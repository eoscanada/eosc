// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/rex"
	"github.com/spf13/cobra"
)

var rexSell = &cobra.Command{
	Use:   "sell [account] [quantity]",
	Short: "Sell REX tokens for EOS tokens.",
	Long:  "Sell REX tokens for EOS tokens. If you have an open `sell` order, this amount will be added to the previous amount to create a single order.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		account := toAccount(args[0], "account")
		quantity := toREXAsset(args[1], "quantity")

		pushEOSCActions(getAPI(), rex.NewSellREX(
			account,
			quantity,
		))
	},
}

func init() {
	rexCmd.AddCommand(rexSell)
}
