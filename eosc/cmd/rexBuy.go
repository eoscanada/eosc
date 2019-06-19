// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/rex"
	"github.com/spf13/cobra"
)

var rexBuy = &cobra.Command{
	Use:   "buy [account] [quantity]",
	Short: "Buy REX tokens using EOS tokens.",
	Long:  "Buy REX tokens using EOS tokens within your REX fund.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		account := toAccount(args[0], "account")
		quantity := toCoreAsset(args[1], "quantity")

		pushEOSCActions(getAPI(), rex.NewBuyREX(
			account,
			quantity,
		))
	},
}

func init() {
	rexCmd.AddCommand(rexBuy)
}
