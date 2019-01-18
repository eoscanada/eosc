// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"fmt"

	"encoding/json"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getScheduledTransactionsCmd = &cobra.Command{
	Use:   "scheduled-transactions",
	Short: "Get scheduled transactions pending for execution.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		lowerBound := viper.GetString("get-scheduled-transactions-cmd-lower_bound")
		limit := viper.GetInt32("get-scheduled-transactions-cmd-limit")

		txs, err := api.GetScheduledTransactionsWithBounds(lowerBound, uint32(limit))
		errorCheck("get scheduled transactions", err)

		data, err := json.MarshalIndent(txs, "", "  ")
		errorCheck("json marshaling", err)

		fmt.Println(string(data))
	},
}

func init() {
	getCmd.AddCommand(getScheduledTransactionsCmd)

	getScheduledTransactionsCmd.Flags().String("lower_bound", "", "The lower bound of the range. Can be a trx_id OR timestamp (as YYYY-MM-DDTHH:MM:SS).")
	getScheduledTransactionsCmd.Flags().Uint32("limit", 100, "The maximum number of deferred transactions to extract at once.")
}
