// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/rex"
	"github.com/spf13/cobra"
)

var rexDefundCPU = &cobra.Command{
	Use:   "defund-cpu [account] [loan number] [quantity]",
	Short: "Remove EOS tokens set for renewal of a CPU loan.",
	Long:  "Remove EOS tokens set for renewal of a CPU loan.",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		account := toAccount(args[0], "account")
		loanNumber := toUint64(args[1], "loan number")
		quantity := toCoreAsset(args[2], "quantity")

		pushEOSCActions(getAPI(), rex.NewDefundCPULoan(
			account,
			loanNumber,
			quantity,
		))
	},
}

func init() {
	rexCmd.AddCommand(rexDefundCPU)
}
