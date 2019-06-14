// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/token"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var transferCmd = &cobra.Command{
	Use:   "transfer [from] [to] [amount]",
	Short: "Transfer from tokens from an account to another",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		contract := toAccount(viper.GetString("transfer-cmd-contract"), "--contract")

		from := toAccount(args[0], "from")
		to := toAccount(args[1], "to")
		quantity := toAssetWithDefaultCoreSymbol(args[2], "quantity")
		memo := viper.GetString("transfer-cmd-memo")

		action := token.NewTransfer(from, to, quantity, memo)
		action.Account = contract

		pushEOSCActions(getAPI(), action)
	},
}

func init() {
	RootCmd.AddCommand(transferCmd)

	transferCmd.Flags().StringP("memo", "m", "", "Memo to attach to the transfer.")
	transferCmd.Flags().StringP("contract", "", "eosio.token", "Contract to send the transfer through. eosio.token is the contract dealing with the native EOS token.")
}
