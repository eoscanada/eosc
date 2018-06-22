package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
)

var getAccountCmd = &cobra.Command{
	Use:   "account [account name]",
	Short: "retrieve account information for a given name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		account, err := api.GetAccount(eos.AccountName(args[0]))
		errorCheck("get account", err)

		data, err := json.MarshalIndent(account, "", "  ")
		errorCheck("json marshal", err)

		fmt.Println(string(data))
	},
}
