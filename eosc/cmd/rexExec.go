// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/rex"
	"github.com/spf13/cobra"
)

var rexExec = &cobra.Command{
	Use:   "exec [account] [max count]",
	Short: "Perform maintenance on the REX contract.",
	Long:  "Perform maintenance on the REX contract (process expired loans or pending sell orders). [max count] needs to be low enough to allow the transaction to be executed within a block.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		account := toAccount(args[0], "account")
		maxCount := toUint16(args[1], "max count")

		pushEOSCActions(getAPI(), rex.NewREXExec(
			account,
			maxCount,
		))
	},
}

func init() {
	rexCmd.AddCommand(rexExec)
}
