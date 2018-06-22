package cmd

import (
	"fmt"
	"os"

	"github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getBalanceCmd = &cobra.Command{
	Use:   "balance [account]",
	Short: "Retrieve currency balance for an account",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		account := eos.AccountName(args[0])
		symbol := viper.GetString("get-balance-cmd-symbol")
		tokenAccount := eos.AccountName(viper.GetString("get-balance-cmd-token-account"))

		balances, err := api.GetCurrencyBalance(account, symbol, tokenAccount)
		if err != nil {
			fmt.Printf("Error: get balance: %s\n", err)
			os.Exit(1)
		}

		for _, asset := range balances {
			fmt.Printf("%s\n", asset)
		}
	},
}

func init() {
	getCmd.AddCommand(getBalanceCmd)

	getBalanceCmd.Flags().StringP("token-account", "", "eosio.token", "Token account")
	getBalanceCmd.Flags().StringP("symbol", "", "", "Only query this symbol. Try EOS")

	for _, flag := range []string{"token-account", "symbol"} {
		if err := viper.BindPFlag("get-balance-cmd-"+flag, getBalanceCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}
