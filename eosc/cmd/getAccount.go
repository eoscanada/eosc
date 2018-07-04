package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var getAccountCmd = &cobra.Command{
	Use:   "account [account name]",
	Short: "retrieve account information for a given name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		accountName := toAccount(args[0], "account name")
		account, err := api.GetAccount(accountName)
		errorCheck("get account", err)

		data, err := json.MarshalIndent(account, "", "  ")
		errorCheck("json marshal", err)

		// TODO: properly display all account details, and fetch a few
		// other things also, to make it a complete picture (like all
		// token balances on eosio.token ?), other tokens ?
		fmt.Println(string(data))
	},
}

func init() {
	getCmd.AddCommand(getAccountCmd)
}
