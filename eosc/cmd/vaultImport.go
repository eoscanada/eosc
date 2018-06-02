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

var vaultImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a private keys to vault",
	Long: `Import a private keys to vault

A vault contains encrypted private keys, and with 'eosc', can be used to
securely sign transactions.

`,
	Run: func(cmd *cobra.Command, args []string) {

		walletFile := viper.GetString("vault-file")
		vault := &eosvault.Vault{}
		var passphrase string
		if _, err := os.Stat(walletFile); err == nil {

			fmt.Println("Loading existing vault from file: ", walletFile)
			vault, err = eosvault.NewVaultFromWalletFile(walletFile)
			if err != nil {
				fmt.Println("ERROR: loading vault from file, ", err)
				os.Exit(1)
			}

			passphrase, err = eosvault.GetPassword("Enter passphrase to unlock vault: ")
			if err != nil {
				fmt.Printf("ERROR: reading passphrase: %s", err)
				os.Exit(1)
			}
			switch vault.SecretBoxWrap {
			case "passphrase":
				err = vault.OpenWithPassphrase(passphrase)
				if err != nil {
					fmt.Printf("ERROR: reading passphrase: %s", err)
					os.Exit(1)
				}
			default:
				fmt.Printf("ERROR: unsupported secretbox wrapping method: %q", vault.SecretBoxWrap)
				os.Exit(1)
			}

			vault.PrintPublicKeys()

		} else {
			fmt.Println("Vault file not found, creating a new wallet")
			vault = eosvault.NewVault()
		}

		if comment := viper.GetString("comment"); comment != "" {
			vault.Comment = comment
		}

		privateKeys, err := capturePrivateKeys()
		if err != nil {
			fmt.Println("ERROR: entering private key:", err)
			os.Exit(1)
		}

		// Got the new keys now

		var newKeys []ecc.PublicKey
		for _, privateKey := range privateKeys {
			vault.AddPrivateKey(privateKey)
			newKeys = append(newKeys, privateKey.PublicKey())
		}

		fmt.Println("Keys imported. Let's secure them before showing the public keys.")

		if passphrase == "" {
			passphrase, err = eosvault.GetPassword("Enter passphrase to encrypt your keys: ")
			if err != nil {
				fmt.Println("ERROR reading password:", err)
				os.Exit(1)
			}

			passphraseConfirm, err := eosvault.GetPassword("Confirm passphrase: ")
			if err != nil {
				fmt.Println("ERROR reading confirmation password:", err)
				os.Exit(1)
			}

			if passphrase != passphraseConfirm {
				fmt.Println("ERROR: passphrase mismatch!")
				os.Exit(1)
			}

			// TODO: make this thing loop.. instead of restarting the whole process..
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

		fmt.Printf("Wallet file %q written. These public keys were ADDED:\n", walletFile)
		for _, pub := range newKeys {
			fmt.Printf("- %s\n", pub.String())
		}
		fmt.Printf("Total keys stored: %d\n", len(vault.KeyBag.Keys))
	},
}

func init() {
	vaultCmd.AddCommand(vaultImportCmd)
	vaultImportCmd.Flags().StringP("comment", "c", "", "Label or comment about this key vault")

	for _, flag := range []string{"comment"} {
		if err := viper.BindPFlag(flag, vaultImportCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}

}

func capturePrivateKeys() ([]*ecc.PrivateKey, error) {
	privateKeys, err := capturePrivateKey(true)
	if err != nil {
		return privateKeys, fmt.Errorf("keys capture, %s", err.Error())
	}
	return privateKeys, nil

}
func capturePrivateKey(isFirst bool) (privateKeys []*ecc.PrivateKey, err error) {
	prompt := "Type your first private key: "
	if !isFirst {
		prompt = "Type your next private key or hit ENTER if you are done: "
	}

	enteredKey, err := eosvault.GetPassword(prompt)
	if err != nil {
		return privateKeys, fmt.Errorf("get password: %s", err.Error())
	}

	if enteredKey == "" {
		return privateKeys, nil
	}

	key, err := ecc.NewPrivateKey(enteredKey)
	if err != nil {
		return privateKeys, fmt.Errorf("new private key: %s", err.Error())
	}

	privateKeys = append(privateKeys, key)
	nextPrivateKeys, err := capturePrivateKey(false)
	if err != nil {
		return privateKeys, err
	}

	privateKeys = append(privateKeys, nextPrivateKeys...)

	return
}
