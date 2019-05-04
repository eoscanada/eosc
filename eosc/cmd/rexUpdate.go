// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/rex"
	"github.com/spf13/cobra"
)

var rexUpdate = &cobra.Command{
	Use:   "update [account]",
	Short: "Update your voting weight.",
	Long:  "Update your voting weight to now include any tokens to which you may be entitled to from increase in the REX pool.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		account := toAccount(args[0], "account")

		pushEOSCActions(getAPI(), rex.NewUpdateREX(
			account,
		))
	},
}

func init() {
	rexCmd.AddCommand(rexUpdate)
}
