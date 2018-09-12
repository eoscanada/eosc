package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var getABICmd = &cobra.Command{
	Use:   "abi [account name]",
	Short: "retrieve the ABI associated with an account",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		accountName := toAccount(args[0], "account name")
		abi, err := api.GetABI(accountName)
		errorCheck("get ABI", err)

		if abi.AccountName == accountName {
			data, err := json.MarshalIndent(abi, "", "  ")
			errorCheck("json marshal", err)
			fmt.Println(string(data))
		} else {
			errorCheck("get abi", fmt.Errorf("unable to find ABI for account %s", accountName))
		}
	},
}

func init() {
	getCmd.AddCommand(getABICmd)
}
