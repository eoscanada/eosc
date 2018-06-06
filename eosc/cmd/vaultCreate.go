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

	"github.com/eoscanada/eos-go/ecc"
	eosvault "github.com/eoscanada/eosc/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// vaultCreateCmd represents the create command
var vaultCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new encrypted EOS keys vault",
	Long: `Create a new encrypted EOS keys vault.

A vault contains encrypted private keys, and with 'eosc', can be used to
securely sign transactions.

`,
	Run: func(cmd *cobra.Command, args []string) {
		walletFile := viper.GetString("vault-file")

		if _, err := os.Stat(walletFile); err == nil {
			fmt.Printf("Wallet file %q already exists, rename it before running `eosc vault create`.\n", walletFile)
			os.Exit(1)
		}

		vault := eosvault.NewVault()
		vault.Comment = viper.GetString("vaultCreateCmd-comment")

		numKeys := viper.GetInt("vaultCreateCmd-keys")
		var newKeys []ecc.PublicKey
		for i := 0; i < numKeys; i++ {
			pubKey, err := vault.NewKeyPair()
			if err != nil {
				fmt.Println("ERROR: creating new keypair:", err)
				os.Exit(1)
			}

			newKeys = append(newKeys, pubKey)
		}

		fmt.Printf("Created %d keys. Let's secure them before showing the public keys.\n", len(newKeys))

		var boxerType = "passphrase-create"
		if viper.GetBool("kms-gcp") {
			boxerType = "kms-gcp"
		}

		boxer, err := eosvault.SecretBoxerForType(boxerType, viper.GetString("kms-keyring"))
		errorCheck(err)

		err = vault.Seal(boxer)
		errorCheck(err)

		err = vault.WriteToFile(walletFile)
		errorCheck(err)

		fmt.Printf("Wallet file %q created. Here are your public keys:\n", walletFile)
		for _, pub := range newKeys {
			fmt.Printf("- %s\n", pub.String())
		}
	},
}

func init() {
	vaultCmd.AddCommand(vaultCreateCmd)

	vaultCreateCmd.Flags().IntP("keys", "k", 1, "Number of keypairs to create")
	vaultCreateCmd.Flags().StringP("comment", "c", "", "Label or comment about this key vault")

	for _, flag := range []string{"keys", "comment"} {
		if err := viper.BindPFlag("vaultCreateCmd-"+flag, vaultCreateCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}
