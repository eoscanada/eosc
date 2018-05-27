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

	"github.com/dgiagio/getpass"
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
		// TODO: check if the viper.GetString("vault-file") exists.
		//   WARN and quit, that's it.
		walletFile := viper.GetString("vault-file")
		if _, err := os.Stat(walletFile); err == nil {
			fmt.Printf("Wallet file %q already exists, rename it before running `eosc vault create`.\n", walletFile)
			os.Exit(1)
		}

		vault := eosvault.NewVault()
		vault.Comment = viper.GetString("comment")

		numKeys := viper.GetInt("keys")
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

		passphrase, err := getpass.GetPassword("Enter passphrase to encrypt your keys: ")
		if err != nil {
			fmt.Println("ERROR reading password:", err)
			os.Exit(1)
		}

		passphraseConfirm, err := getpass.GetPassword("Confirm passphrase: ")
		if err != nil {
			fmt.Println("ERROR reading confirmation password:", err)
			os.Exit(1)
		}

		if passphrase != passphraseConfirm {
			fmt.Println("ERROR: passphrase mismatch!")
			os.Exit(1)
		}

		err = vault.SealWithPassphrase(passphrase)
		if err != nil {
			fmt.Println("ERROR sealing keys:", err)
			os.Exit(1)
		}

		err = vault.WriteToFile(walletFile)
		if err != nil {
			fmt.Printf("ERROR writing to file %q: %s\n", walletFile, err)
			os.Exit(1)
		}

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
		if err := viper.BindPFlag(flag, vaultCreateCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}
