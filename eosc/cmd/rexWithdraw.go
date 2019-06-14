// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/rex"
	"github.com/spf13/cobra"
)

var rexWithdraw = &cobra.Command{
	Use:   "withdraw [account] [quantity]",
	Short: "Withdraw EOS tokens from your REX fund.",
	Long:  "Withdraw EOS tokens from your REX fund to your liquid EOS balance.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		account := toAccount(args[0], "account")
		quantity := toCoreAsset(args[1], "quantity")

		pushEOSCActions(getAPI(), rex.NewWithdraw(
			account,
			quantity,
		))
	},
}

func init() {
	rexCmd.AddCommand(rexWithdraw)
}
