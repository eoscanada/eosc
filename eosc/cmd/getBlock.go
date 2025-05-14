// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var getBlockCmd = &cobra.Command{
	Use:   "block [block id | block height]",
	Short: "Get block data at a given height, or directly with a block hash",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		block, err := api.GetBlockByNumOrIDRaw(context.Background(), args[0])
		errorCheck("get block", err)

		data, err := json.MarshalIndent(block, "", "  ")
		errorCheck("json marshaling", err)

		fmt.Println(string(data))
	},
}

func init() {
	getCmd.AddCommand(getBlockCmd)

	// getBlockCmd.Flags().BoolP("json", "", false, "return producers info in json")
}
