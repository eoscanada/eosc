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
	"github.com/spf13/viper"
)

var getBlockCmd = &cobra.Command{
	Use:   "get [block id | block height]",
	Short: "get block info",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		block, err := api.GetBlockByNumOrID(args[0])
		errorCheck("get block", err)

		if viper.GetBool("getBlockCmd.json") {
			data, err := json.MarshalIndent(block, "", "  ")
			errorCheck("json marshal", err)

			fmt.Println(string(data))
		} else {
			fmt.Printf("Block [%d] produced by [%s] as [%s]\n", block.BlockHeader.BlockNumber(), block.BlockHeader.Producer, block.Timestamp)
		}
	},
}

var blockCmd = &cobra.Command{
	Use:   "block",
	Short: "block related commands",
	Long:  `block related commands`,
}

func init() {
	// RootCmd.AddCommand(blockCmd)
	// blockCmd.AddCommand(getBlockCmd)

	getBlockCmd.Flags().BoolP("json", "", false, "return producers info in json")

	for _, flag := range []string{"json", "limit"} {
		if err := viper.BindPFlag("getBlockCmd."+flag, getBlockCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}

}
