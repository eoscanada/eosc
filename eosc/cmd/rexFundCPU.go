// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/rex"
	"github.com/spf13/cobra"
)

var rexFundCPU = &cobra.Command{
	Use:   "fund-cpu [account] [loan number] [quantity]",
	Short: "Set EOS tokens to renew a CPU loan upon expiry.",
	Long:  "Set an amount of EOS tokens from your REX fund to be used to renew a CPU loan upon expiry.",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		account := toAccount(args[0], "account")
		loanNumber := toUint64(args[1], "loan number")
		quantity := toCoreAsset(args[2], "quantity")

		pushEOSCActions(getAPI(), rex.NewFundCPULoan(
			account,
			loanNumber,
			quantity,
		))
	},
}

func init() {
	rexCmd.AddCommand(rexFundCPU)
}
