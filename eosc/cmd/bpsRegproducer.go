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

	"os"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var bpsRegproducerCmd = &cobra.Command{
	Use:   "register [account_name] [public_key] [website_url]",
	Short: "Register account as a block producer",
	Long:  `Register account as a block producer`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		api, err := api()
		if err != nil {
			fmt.Printf("Error initiating api, %s\n", err.Error())
			os.Exit(1)
		}

		accountName := eos.AccountName(args[0])
		publicKey, err := ecc.NewPublicKey(args[1])
		if err != nil {
			fmt.Printf("Error public key, %s\n", err.Error())
			os.Exit(1)
		}
		websiteURL := args[2]

		// Validate the input params

		_, err = api.SignPushActions(
			system.NewRegProducer(accountName, publicKey, websiteURL),
		)

		if err != nil {
			fmt.Printf("Producer registration , %s\n", err.Error())
			os.Exit(1)
		}

		fmt.Println("Producer registration sent to chain")
	},
}

func init() {
	bpsCmd.AddCommand(bpsRegproducerCmd)
}
