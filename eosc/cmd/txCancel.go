// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var txCancelCmd = &cobra.Command{
	Use:   `cancel [cancelling_authority] [transaction_id]`,
	Short: "Cancels a delayed transaction.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		authority := toPermissionLevel(args[0], "cancelling_authority")
		transactionID := toSHA256Bytes(args[1], "transaction_id")

		api := getAPI()
		pushEOSCActions(api, system.NewCancelDelay(authority, transactionID))
	},
}

func init() {
	txCmd.AddCommand(txCancelCmd)
}
