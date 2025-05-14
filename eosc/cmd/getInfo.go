// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var getInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Retrieve blockchain infos, like head block, chain ID, etc..",
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		info, err := api.GetInfo(context.Background())
		errorCheck("get info", err)

		data, err := json.MarshalIndent(info, "", "  ")
		errorCheck("json marshal", err)

		fmt.Println(string(data))
	},
}

func init() {
	getCmd.AddCommand(getInfoCmd)
}
