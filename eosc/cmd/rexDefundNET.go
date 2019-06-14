// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/rex"
	"github.com/spf13/cobra"
)

var rexDefundNet = &cobra.Command{
	Use:   "defund-net [account] [loan number] [quantity]",
	Short: "Remove EOS tokens set for renewal of a Network loan.",
	Long:  "Remove EOS tokens set for renewal of a Network loan.",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		account := toAccount(args[0], "account")
		loanNumber := toUint64(args[1], "loan number")
		quantity := toCoreAsset(args[2], "quantity")

		pushEOSCActions(getAPI(), rex.NewDefundNetLoan(
			account,
			loanNumber,
			quantity,
		))
	},
}

func init() {
	rexCmd.AddCommand(rexDefundNet)
}
