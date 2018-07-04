// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/token"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var transferCmd = &cobra.Command{
	Use:   "transfer [from] [to] [amount]",
	Short: "Transfer from tokens from an account to another",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {

		from := toAccount(args[0], "from")
		to := toAccount(args[1], "to")
		quantity, err := eos.NewEOSAssetFromString(args[2])
		errorCheck("invalid amount", err)
		memo := viper.GetString("transfer-cmd-memo")

		api := apiWithWallet()

		action := token.NewTransfer(from, to, quantity, memo)
		action.Account = toAccount(viper.GetString("transfer-cmd-contract"), "--contract")
		pushEOSCActions(api, action)
	},
}

func init() {
	RootCmd.AddCommand(transferCmd)

	transferCmd.Flags().StringP("memo", "m", "", "Memo to attach to the transfer.")
	transferCmd.Flags().StringP("contract", "", "eosio.token", "Contract to send the transfer through. eosio.token is the contract dealing with the native EOS token.")

	for _, flag := range []string{"memo", "contract"} {
		if err := viper.BindPFlag("transfer-cmd-"+flag, forumPostCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}
