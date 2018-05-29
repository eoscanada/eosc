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

var bpsListProducerCmd = &cobra.Command{
	Use:   "list",
	Short: "List the producers",
	Long:  `List the producers`,
	Run: func(cmd *cobra.Command, args []string) {
		api, err := api()
		if err != nil {
			fmt.Printf("Error initiating api, %s\n", err.Error())
			os.Exit(1)
		}

		response, err := api.GetTableRows(
			eos.GetTableRowsRequest{
				Scope: "eosio",
				Code:  "eosio",
				Table: "producers",
				JSON:  true,
				Limit: uint32(viper.GetInt("limit")),
			},
		)

		if err != nil {
			fmt.Printf("Get table rows , %s\n", err.Error())
			os.Exit(1)
		}

		if viper.GetBool("json") {
			data, err := json.MarshalIndent(response.Rows, "", "    ")
			if err != nil {
				fmt.Printf("JSON generation , %s\n", err.Error())
				os.Exit(1)
			}
			fmt.Println(string(data))
		} else {
			var producers []interface{}
			if err := json.Unmarshal(response.Rows, &producers); err != nil {
				fmt.Printf("JSON unmarshall , %s\n", err.Error())
				os.Exit(1)
			}

			for _, p := range producers {
				producerMap := p.(map[string]interface{})
				fmt.Printf("Producer [%s]: total vote [%s]\n", producerMap["owner"], producerMap["total_votes"])
			}

		}
	},
}

func init() {
	bpsCmd.AddCommand(bpsListProducerCmd)

	bpsListProducerCmd.Flags().BoolP("json", "", false, "return producers info in json")
	bpsListProducerCmd.Flags().IntP("limit", "", 50, "maximum producers that will be return")

	for _, flag := range []string{"json", "limit"} {
		if err := viper.BindPFlag(flag, bpsListProducerCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}

}
