// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/rex"
	"github.com/spf13/cobra"
)

var rexDeposit = &cobra.Command{
	Use:   "deposit [account] [quantity]",
	Short: "Deposit EOS tokens into your REX fund.",
	Long:  "Deposit EOS tokens into your REX fund, to be used to purchase REX tokens.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		account := toAccount(args[0], "account")
		quantity := toCoreAsset(args[1], "quantity")

		pushEOSCActions(getAPI(), rex.NewDeposit(
			account,
			quantity,
		))
	},
}

func init() {
	rexCmd.AddCommand(rexDeposit)
}
