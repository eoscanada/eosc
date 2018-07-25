// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getTableCmd = &cobra.Command{
	Use:   "table [contract] [scope] [table]",
	Short: "Fetch data from a table in a contract on chain",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		response, err := api.GetTableRows(
			eos.GetTableRowsRequest{
				Code:  args[0],
				Scope: args[1],
				Table: args[2],
				JSON:  true,
				Limit: uint32(viper.GetInt("get-table-cmd-limit")),
			},
		)
		errorCheck("get table rows", err)

		data, err := json.MarshalIndent(response, "", "  ")
		errorCheck("json marshal", err)

		fmt.Println(string(data))
	},
}

func init() {
	getCmd.AddCommand(getTableCmd)

	getTableCmd.Flags().IntP("limit", "", 100, "Maximum number of rows to return.")

	for _, flag := range []string{"limit"} {
		if err := viper.BindPFlag("get-table-cmd-"+flag, getTableCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}

}
