// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"encoding/json"
	"fmt"

	"sort"

	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var voteListProducersCmd = &cobra.Command{
	Use:   "list-producers",
	Short: "Retrieve the list of registered producers.",
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

	producers, err := getProducersTable(api)
	errorCheck("get producers table", err)

	if viper.GetBool("vote-list-cmd-json") {
		data, err := json.MarshalIndent(producers, "", "  ")
		errorCheck("json marshal", err)

		fmt.Println(string(data))
	} else {
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
	voteCmd.AddCommand(voteListProducersCmd)

	voteListProducersCmd.Flags().BoolP("sort", "s", false, "sort producers")
	voteListProducersCmd.Flags().BoolP("json", "j", false, "return producers info in json")
}
