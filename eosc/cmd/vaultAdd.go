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
	"github.com/eoscanada/eosc/cli"
	eosvault "github.com/eoscanada/eosc/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var vaultAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add private keys to vault from command line input.",
	Run: func(cmd *cobra.Command, args []string) {

		walletFile := viper.GetString("global-vault-file")

		fmt.Println("Loading existing vault from file:", walletFile)
		vault, err := eosvault.NewVaultFromWalletFile(walletFile)
		errorCheck("loading vault from file", err)

		boxer, err := eosvault.SecretBoxerForType(vault.SecretBoxWrap, viper.GetString("vault-cmd-kms-gcp-keypath"))
		errorCheck("missing parameters", err)

		err = vault.Open(boxer)
		errorCheck("opening vault", err)

		vault.PrintPublicKeys()

		privateKeys, err := capturePrivateKeys()
		errorCheck("entering private keys", err)

		var newKeys []ecc.PublicKey
		for _, privateKey := range privateKeys {
			vault.AddPrivateKey(privateKey)
			newKeys = append(newKeys, privateKey.PublicKey())
		}

		err = vault.Seal(boxer)
		errorCheck("sealing vault", err)

		err = vault.WriteToFile(walletFile)
		errorCheck("writing vault file", err)

		fmt.Printf("Wallet file %q written. These keys were ADDED:\n", walletFile)
		for _, pub := range newKeys {
			fmt.Printf("- %s\n", pub.String())
		}
		fmt.Printf("Total keys stored: %d\n", len(vault.KeyBag.Keys))
	},
}

func init() {
	vaultCmd.AddCommand(vaultAddCmd)
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

	enteredKey, err := cli.GetPassword(prompt)
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
