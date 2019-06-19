// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/rex"
	"github.com/spf13/cobra"
)

var rexRentNet = &cobra.Command{
	Use:   "rent-net [payer] [receiver] [quantity] [loan fund]",
	Short: "Rent Network for the account [receiver].",
	Long:  "Rent Network for the account [receiver]. If you don't want to set up an automatic renewal upon expiry, enter 0 for [loan fund]. Otherwise you can set up an amount now, or use `eosc rex fundnetloan` to provide tokens to renew the loan.",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		payer := toAccount(args[0], "payer")
		receiver := toAccount(args[1], "receiver")
		quantity := toCoreAsset(args[2], "quantity")
		loanFund := toCoreAsset(args[3], "loan fund")

		pushEOSCActions(getAPI(), rex.NewRentNet(
			payer,
			receiver,
			quantity,
			loanFund,
		))
	},
}

func init() {
	rexCmd.AddCommand(rexRentNet)
}
