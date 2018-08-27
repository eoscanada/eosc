// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/spf13/cobra"
)

// txCmd represents the tx command
var txCmd = &cobra.Command{
	Use:   "tx",
	Short: "Transactions-related commands, like signing, pushing, reading, etc..",
}

func init() {
	RootCmd.AddCommand(txCmd)

	//txCmd.PersistentFlags().StringP("target-contract", "", "eostxdapp", "Target account hosting the eosio.tx code")

	// for _, flag := range []string{"target-contract"} {
	// 	if err := viper.BindPFlag("tx-cmd-"+flag, txCmd.PersistentFlags().Lookup(flag)); err != nil {
	// 		panic(err)
	// 	}
	// }
}
