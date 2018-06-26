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
	"github.com/spf13/viper"
)

var systemNewAccountCmd = &cobra.Command{
	Use:   "newaccount [creator] [new_account_name]",
	Short: "Create a new account",
	Long: `Create a new account

Specify the authority structure with either '--auth-file' or '--auth-key'.

With --auth-key, the provided EOS public key will be used for both the
owner and active permissions.

With --auth-file, you can create authority structures for both owner
and active, from the start. Here is a sample auth file in YAML:

---
owner:
  threshold: 2
  keys:
  - key: EOS6MRyAjQq8ud7hVNYcfn................tHuGYqET5GDW5CV
    weight: 1
  waits:
  - wait_sec: 300
    weight: 1
active:
  threshold: 1
  accounts:
  - permission:
      actor: otheraccount
      permission: active
    weight: 1
---

`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		creator := eos.AccountName(args[0])
		newAccount := eos.AccountName(args[1])

		var actions []*eos.Action
		authFile := viper.GetString("system-newaccount-cmd-auth-file")
		authKey := viper.GetString("system-newaccount-cmd-auth-key")
		if authKey == "" && authFile == "" {
			fmt.Println("Error: pass one of --auth-file or --auth-key")
			os.Exit(1)
		}

		if authKey != "" && authFile != "" {
			fmt.Println("Error: pass either --auth-file or --auth-key")
			os.Exit(1)
		}

		if authFile != "" {
			// load from YAML
			var authStruct struct {
				Owner  eos.Authority `json:"owner"`
				Active eos.Authority `json:"active"`
			}
			err := loadYAMLOrJSONFile(authFile, &authStruct)
			errorCheck("auth-file invalid", err)

			if authStruct.Owner.Threshold == 0 {
				errorCheck("auth-file invalid", fmt.Errorf("owner struct missing?"))
			}

			if authStruct.Active.Threshold == 0 {
				errorCheck("auth-file invalid", fmt.Errorf("active struct missing?"))
			}

			actions = append(actions, system.NewCustomNewAccount(creator, newAccount, authStruct.Owner, authStruct.Active))
		} else {
			// authKey then
			pubKey, err := ecc.NewPublicKey(authKey)
			errorCheck("parsing public key", err)

			actions = append(actions, system.NewNewAccount(creator, newAccount, pubKey))
		}

		cpuStakeStr := viper.GetString("system-newaccount-cmd-stake-cpu")
		netStakeStr := viper.GetString("system-newaccount-cmd-stake-net")

		if cpuStakeStr == "" {
			errorCheck("missing argument", fmt.Errorf("--stake-cpu missing"))
		}
		if netStakeStr == "" {
			errorCheck("missing argument", fmt.Errorf("--stake-net missing"))
		}

		cpuStake, err := eos.NewEOSAssetFromString(cpuStakeStr)
		errorCheck("--stake-cpu invalid", err)
		netStake, err := eos.NewEOSAssetFromString(netStakeStr)
		errorCheck("--stake-net invalid", err)

		doTransfer := viper.GetBool("system-newaccount-cmd-transfer")
		actions = append(actions, system.NewDelegateBW(creator, newAccount, cpuStake, netStake, doTransfer))

		buyRAM := viper.GetString("system-newaccount-cmd-buy-ram")
		if buyRAM != "" {
			buyRAMAmount, err := eos.NewEOSAssetFromString(buyRAM)
			errorCheck("--buy-ram invalid", err)

			actions = append(actions, system.NewBuyRAM(creator, newAccount, uint64(buyRAMAmount.Amount)))
		} else {
			buyRAMBytes := viper.GetInt("system-newaccount-cmd-buy-ram-kbytes")
			actions = append(actions, system.NewBuyRAMBytes(creator, newAccount, uint32(buyRAMBytes*1024)))
		}

		api := apiWithWallet()

		pushEOSCActions(api, actions...)
	},
}

func init() {
	systemCmd.AddCommand(systemNewAccountCmd)

	systemNewAccountCmd.Flags().StringP("auth-file", "", "", "File containing owner and active permissions authorities. See example above.")
	systemNewAccountCmd.Flags().StringP("auth-key", "", "", "Public key to use for both owner and active permissions.")
	systemNewAccountCmd.Flags().StringP("stake-cpu", "", "", "Amount of EOS to stake for CPU bandwidth (required)")
	systemNewAccountCmd.Flags().StringP("stake-net", "", "", "Amount of EOS to stake for Network bandwidth (required)")
	systemNewAccountCmd.Flags().IntP("buy-ram-kbytes", "", 8, "The amount of RAM kibibytes (KiB) to purchase for the new account.  Defaults to 8 KiB.")
	systemNewAccountCmd.Flags().StringP("buy-ram", "", "", "The amount of EOS to spend to buy RAM for the new account (at current EOS/RAM market price)")
	systemNewAccountCmd.Flags().BoolP("transfer", "", false, "Transfer voting power and right to unstake EOS to receiver")

	for _, flag := range []string{"stake-cpu", "stake-net", "buy-ram-kbytes", "buy-ram", "transfer", "auth-file", "auth-key"} {
		if err := viper.BindPFlag("system-newaccount-cmd-"+flag, systemNewAccountCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}

}
