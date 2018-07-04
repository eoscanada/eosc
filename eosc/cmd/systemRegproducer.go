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
	"fmt"

	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var systemRegisterProducerCmd = &cobra.Command{
	Use:   "register [account_name] [public_key] [website_url]",
	Short: "Register account as a block producer",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		api := apiWithWallet()

		accountName := toAccount(args[0], "account name")
		publicKey, err := ecc.NewPublicKey(args[1])
		errorCheck(fmt.Sprintf("%q invalid public key", args[1]), err)
		websiteURL := args[2]

		pushEOSCActions(api,
			system.NewRegProducer(accountName, publicKey, websiteURL),
		)
	},
}

func init() {
	systemCmd.AddCommand(systemRegisterProducerCmd)
}
