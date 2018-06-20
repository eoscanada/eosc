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

	"sort"

	"strconv"

	"github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var voteListProducerCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve registered producers",
	Run:   run,
}

type producers []map[string]interface{}

func (p producers) Len() int      { return len(p) }
func (p producers) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p producers) Less(i, j int) bool {
	iv, _ := strconv.ParseFloat(p[i]["total_votes"].(string), 64)
	jv, _ := strconv.ParseFloat(p[j]["total_votes"].(string), 64)
	return iv > jv
}

var run = func(cmd *cobra.Command, args []string) {
	api := getAPI()

	response, err := api.GetTableRows(
		eos.GetTableRowsRequest{
			Scope: "eosio",
			Code:  "eosio",
			Table: "producers",
			JSON:  true,
			Limit: 5000,
		},
	)

	if err != nil {
		fmt.Printf("Get table rows , %s\n", err.Error())
		os.Exit(1)
	}

	if viper.GetBool("vote-list-cmd-json") {
		data, err := json.MarshalIndent(response.Rows, "", "    ")
		if err != nil {
			fmt.Println("Error: json marshalling:", err)
			os.Exit(1)
		}
		fmt.Println(string(data))
	} else {
		var producers producers
		if err := json.Unmarshal(response.Rows, &producers); err != nil {
			fmt.Println("Error: json marshalling:", err)
			os.Exit(1)
		}

		if viper.GetBool("vote-list-cmd-sort") {
			sort.Slice(producers, producers.Less)
		}

		fmt.Println("List of producers registered to receive votes:")
		for _, p := range producers {
			fmt.Printf("- %s (key: %s)\n", p["owner"], p["producer_key"])
		}
		fmt.Printf("Total of %d registered producers\n", len(producers))

	}
}

func init() {
	voteCmd.AddCommand(voteListProducerCmd)

	voteListProducerCmd.Flags().BoolP("sort", "s", false, "sort producers")
	voteListProducerCmd.Flags().BoolP("json", "j", false, "return producers info in json")

	for _, flag := range []string{"json", "sort"} {
		if err := viper.BindPFlag("vote-list-cmd-"+flag, voteListProducerCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}
