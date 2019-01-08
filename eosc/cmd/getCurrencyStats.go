package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/eoscanada/eosc/cli"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getCurrencyStatsCmd = &cobra.Command{
	Use:   "currency-stats [account name] [symbol]",
	Short: "retrieve currency statistics",
	Long:  "retrieve currency statistics (supply, max. supply, issuer name).",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		accountName := toAccount(args[0], "account name")
		stats, err := api.GetCurrencyStats(accountName, args[1])
		errorCheck("get currency-stats", err)

		if stats == nil {
			fmt.Println("Currency not found.")
			return
		}

		if viper.GetBool("get-currency-stats-cmd-json") == true {
			data, err := json.MarshalIndent(stats, "", "  ")
			errorCheck("json marshal", err)
			fmt.Println(string(data))
			return
		}

		cfg := columnize.DefaultConfig()
		fmt.Println(cli.FormatCurrencyStats(stats, cfg))
	},
}

func init() {
	getCmd.AddCommand(getCurrencyStatsCmd)
	getCurrencyStatsCmd.Flags().BoolP("json", "", false, "pass if you wish to see account printed as json")

	for _, flag := range []string{"json"} {
		if err := viper.BindPFlag("get-currency-stats-cmd-"+flag, getCurrencyStatsCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}
