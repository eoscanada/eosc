package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "account related commands",
	Long:  `account related commands`,
}

var getAccountCmd = &cobra.Command{
	Use:   "get",
	Short: "retrieve account information",
	Long:  `retrieve account information`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api, err := api()
		api.Debug = true
		if err != nil {
			fmt.Printf("Error initiating api, %s\n", err.Error())
			os.Exit(1)
		}

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

func init() {
	RootCmd.AddCommand(accountCmd)
	accountCmd.AddCommand(getAccountCmd)
}
