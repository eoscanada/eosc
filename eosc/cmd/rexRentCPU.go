// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/rex"
	"github.com/spf13/cobra"
)

var rexRentCPU = &cobra.Command{
	Use:   "rent-cpu [payer] [receiver] [quantity] [loan fund]",
	Short: "Rent CPU for the account [receiver].",
	Long:  "Rent CPU for the account [receiver]. If you don't want to set up an automatic renewal upon expiry, enter 0 for [loan fund]. Otherwise you can set up an amount now, or use `eosc rex fund-cpu` to provide tokens to renew the loan.",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		payer := toAccount(args[0], "payer")
		receiver := toAccount(args[1], "receiver")
		quantity := toCoreAsset(args[2], "quantity")
		loanFund := toCoreAsset(args[3], "loan fund")

		pushEOSCActions(getAPI(), rex.NewRentCPU(
			payer,
			receiver,
			quantity,
			loanFund,
		))
	},
}

func init() {
	rexCmd.AddCommand(rexRentCPU)
}
