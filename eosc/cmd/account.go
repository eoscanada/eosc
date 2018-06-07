package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
)

var getAccountCmd = &cobra.Command{
	Use:   "account [account name]",
	Short: "retrieve account information for a given name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := api()

		account, err := api.GetAccount(eos.AccountName(args[0]))
		if err != nil {
			fmt.Printf("Error: get account , %s\n", err.Error())
			os.Exit(1)
		}
		data, err := json.MarshalIndent(account, "", "  ")
		if err != nil {
			fmt.Printf("Error: json conversion , %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Println(string(data))
	},
}
