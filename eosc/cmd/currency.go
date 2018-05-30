package cmd

import (
	"fmt"
	"os"

	"math"

	"github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var currencyCmd = &cobra.Command{
	Use:   "currency",
	Short: "currency related commands",
	Long:  `currency related commands`,
}

var getCurrencyBalanceCmd = &cobra.Command{
	Use:   "balance [account]",
	Short: "retrieve currency balance for an account",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		api := api()

		balances, err := api.GetCurrencyBalance(eos.AccountName(args[0]), viper.GetString("currencyCmd.symbol"), eos.AccountName(viper.GetString("currencyCmd.currency-account")))
		if err != nil {
			fmt.Printf("Error: get account , %s\n", err.Error())
			os.Exit(1)
		}

		for _, asset := range balances {
			fmt.Printf("- Balance for %s is %.4f %s\n", args[0], float64(asset.Amount)/float64(math.Pow(10, float64(asset.Precision))), asset.Symbol.Symbol)
		}

		//data, err := json.MarshalIndent(balance, "", "  ")
		//if err != nil {
		//	fmt.Printf("Error: json conversion , %s\n", err.Error())
		//	os.Exit(1)
		//}
		//fmt.Println(string(data))
	},
}

func init() {
	// RootCmd.AddCommand(currencyCmd)
	// currencyCmd.AddCommand(getCurrencyBalanceCmd)

	currencyCmd.Flags().StringP("currency-account", "", "eosio.token", "account owning the currency")
	currencyCmd.Flags().StringP("symbol", "", "EOS", "symbol representing the currency")

	for _, flag := range []string{"currency-account", "symbol"} {
		if err := viper.BindPFlag("currencyCmd."+flag, currencyCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}
