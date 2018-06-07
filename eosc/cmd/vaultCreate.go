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

// vaultCreateCmd represents the create command
var vaultCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new encrypted EOS keys vault",
	Long: `Create a new encrypted EOS keys vault.

A vault contains encrypted private keys, and with 'eosc', can be used to
securely sign transactions.

`,
	Run: func(cmd *cobra.Command, args []string) {
		walletFile := viper.GetString("global-vault-file")

		if _, err := os.Stat(walletFile); err == nil {
			fmt.Printf("Wallet file %q already exists, rename it before running `eosc vault create`.\n", walletFile)
			os.Exit(1)
		}

		var wrapType = viper.GetString("vault-create-cmd-vault-type")
		var boxer eosvault.SecretBoxer

		switch wrapType {
		case "kms-gcp":
			keyring := viper.GetString("kms-gcp-keyring")
			if keyring == "" {
				errorCheck("missing parameter", fmt.Errorf("--kms-gcp-keyring is required with --vault-type=kms-gcp"))
			}
			boxer = eosvault.NewKMSGCPBoxer(keyring)

		case "passphrase":
			password, err := cli.GetEncryptPassphrase()
			errorCheck("password input", err)

			boxer = eosvault.NewPassphraseBoxer(password)

		default:
			fmt.Printf(`Invalid vault type: %q, please use one of: "passphrase", "kms-gcp"\n`, wrapType)
			os.Exit(1)
		}

		vault := eosvault.NewVault()
		vault.Comment = viper.GetString("vault-create-cmd-comment")

		var newKeys []ecc.PublicKey

		doImport := viper.GetBool("vault-create-cmd-import")
		if doImport {
			privateKeys, err := capturePrivateKeys()
			errorCheck("entering private key", err)

			for _, privateKey := range privateKeys {
				vault.AddPrivateKey(privateKey)
				newKeys = append(newKeys, privateKey.PublicKey())
			}

			fmt.Printf("Imported %d keys. Let's secure them before showing the public keys.\n", len(newKeys))

		} else {
			numKeys := viper.GetInt("vault-create-cmd-keys")
			for i := 0; i < numKeys; i++ {
				pubKey, err := vault.NewKeyPair()
				errorCheck("creating new keypair", err)

				newKeys = append(newKeys, pubKey)
			}
			fmt.Printf("Created %d keys. Let's secure them before showing the public keys.\n", len(newKeys))
		}

		errorCheck("sealing vault", vault.Seal(boxer))
		errorCheck("writing wallet file", vault.WriteToFile(walletFile))

		fmt.Printf("Wallet file %q created. Here are your public keys:\n", walletFile)
		for _, pub := range newKeys {
			fmt.Printf("- %s\n", pub.String())
		}
	},
}

func init() {
	vaultCmd.AddCommand(vaultCreateCmd)

	vaultCreateCmd.Flags().IntP("keys", "k", 1, "Number of keypairs to create")
	vaultCreateCmd.Flags().BoolP("import", "i", false, "Whether to import keys instead of creating them. This takes precedence over --keys, and private keys will be inputted on the command line.")
	vaultCreateCmd.Flags().StringP("comment", "c", "", "Comment field in the vault's json file.")
	vaultCreateCmd.Flags().StringP("vault-type", "t", "passphrase", "Vault type. One of: passphrase, kms-gcp")

	for _, flag := range []string{"keys", "comment", "vault-type", "import"} {
		if err := viper.BindPFlag("vault-create-cmd-"+flag, vaultCreateCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}
