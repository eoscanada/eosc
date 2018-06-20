// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"

	"os"

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

		if err != nil {
			fmt.Printf("Error: get table rows: %s\n", err)
			os.Exit(1)
		}

		data, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			fmt.Printf("Error: json marshal: %s\n", err)
			os.Exit(1)
		}

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
