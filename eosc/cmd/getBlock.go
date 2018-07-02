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
	"fmt"

	"encoding/json"

	"github.com/spf13/cobra"
)

var getBlockCmd = &cobra.Command{
	Use:   "block [block id | block height]",
	Short: "Get block data at a given height, or directly with a block hash",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		block, err := api.GetBlockByNumOrIDRaw(args[0])
		errorCheck("get block", err)

		data, err := json.MarshalIndent(block, "", "  ")
		errorCheck("json marshaling", err)

		fmt.Println(string(data))
	},
}

func init() {
	getCmd.AddCommand(getBlockCmd)

	// getBlockCmd.Flags().BoolP("json", "", false, "return producers info in json")

	// for _, flag := range []string{"json"} {
	// 	if err := viper.BindPFlag("get-block-cmd-"+flag, getBlockCmd.Flags().Lookup(flag)); err != nil {
	// 		panic(err)
	// 	}
	// }
}
