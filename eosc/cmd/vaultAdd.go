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
	"github.com/eoscanada/eosc/cli"
	eosvault "github.com/eoscanada/eosc/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var vaultAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add private keys to an existing vault taking input from the shell",
	Run: func(cmd *cobra.Command, args []string) {

		walletFile := viper.GetString("global-vault-file")

		fmt.Println("Loading existing vault from file:", walletFile)
		vault, err := eosvault.NewVaultFromWalletFile(walletFile)
		errorCheck("loading vault from file", err)

		boxer, err := eosvault.SecretBoxerForType(vault.SecretBoxWrap, viper.GetString("global-kms-gcp-keypath"))
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

		vaultWrittenReport(walletFile, newKeys, len(vault.KeyBag.Keys))
	},
}

func init() {
	vaultCmd.AddCommand(vaultAddCmd)
}

func capturePrivateKeys() (out []*ecc.PrivateKey, err error) {
	fmt.Println("")
	fmt.Println("PLEASE READ:")
	fmt.Println("We are now going to ask you to paste your private keys, one at a time.")
	fmt.Println("They will not be shown on screen.")
	fmt.Println("Please verify that the public keys printed on screen correspond to what you have noted")
	fmt.Println("")

	first := true
	for {
		privKey, err := capturePrivateKey(first)
		if err != nil {
			return out, fmt.Errorf("capture privkeys: %s", err)
		}
		first = false

		if privKey == nil {
			return out, nil
		}
		out = append(out, privKey)
	}

	panic("dead code path")
}
func capturePrivateKey(isFirst bool) (privateKey *ecc.PrivateKey, err error) {
	prompt := "Paste your first private key: "
	if !isFirst {
		prompt = "Paste your next private key or hit ENTER if you are done: "
	}

	enteredKey, err := cli.GetPassword(prompt)
	if err != nil {
		return nil, fmt.Errorf("get private key: %s", err)
	}

	if enteredKey == "" {
		return nil, nil
	}

	key, err := ecc.NewPrivateKey(enteredKey)
	if err != nil {
		return nil, fmt.Errorf("import private key: %s", err)
	}

	fmt.Printf("- Scanned private key corresponding to %s\n", key.PublicKey().String())

	return key, nil
}

func vaultWrittenReport(walletFile string, newKeys []ecc.PublicKey, totalKeys int) {
	fmt.Println("")
	fmt.Printf("Wallet file %q written to disk.\n", walletFile)
	fmt.Println("Here are the keys that were ADDED during this operation (use `list` to see them all):")
	for _, pub := range newKeys {
		fmt.Printf("- %s\n", pub.String())
	}

	fmt.Printf("Total keys stored: %d\n", totalKeys)
}
