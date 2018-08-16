// Copyright © 2018 EOS Canada <info@eoscanada.com>

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

You can create a passphrase protected vault with:

    eosc vault create --keys=2

This uses the default --vault-type=passphrase

You can create a Google Cloud Platform KMS-wrapped vault with:

    eosc vault create --keys=2 --vault-type=kms-gcp --kms-gcp-keypath projects/.../locations/.../keyRings/.../cryptoKeys/name

You can then use this vault for the different eosc operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		walletFile := viper.GetString("global-vault-file")

		if _, err := os.Stat(walletFile); err == nil {
			fmt.Printf("Wallet file %q already exists, rename it before running `eosc vault create`.\n", walletFile)
			os.Exit(1)
		}

		var wrapType = viper.GetString("vault-create-cmd-vault-type")
		var boxer eosvault.SecretBoxer

		kmsGCPKeypath := viper.GetString("global-kms-gcp-keypath")
		if wrapType == "kms-gcp" && kmsGCPKeypath == "" {
			errorCheck("missing parameter", fmt.Errorf("--kms-gcp-keypath is required with --vault-type=kms-gcp"))
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

			fmt.Printf("Imported %d keys.\n", len(newKeys))

		} else {
			numKeys := viper.GetInt("vault-create-cmd-keys")

			if numKeys == 0 {
				errorCheck("specify either --keys or --import", fmt.Errorf("create a vault with 0 keys?"))
			}

			for i := 0; i < numKeys; i++ {
				pubKey, err := vault.NewKeyPair()
				errorCheck("creating new keypair", err)

				newKeys = append(newKeys, pubKey)
			}
			fmt.Printf("Created %d keys. They will be shown when encrypted and written to disk successfully.\n", len(newKeys))
		}

		switch wrapType {
		case "kms-gcp":
			fmt.Println("Doing the KMS GCP dance")
			boxer = eosvault.NewKMSGCPBoxer(kmsGCPKeypath)

		case "passphrase":
			fmt.Println("")
			fmt.Println("You will be asked to provide a passphrase to secure your newly created vault.")
			fmt.Println("Make sure you make it long and strong.")
			fmt.Println("")
			password, err := cli.GetEncryptPassphrase()
			errorCheck("password input", err)

			boxer = eosvault.NewPassphraseBoxer(password)

		default:
			fmt.Printf(`Invalid vault type: %q, please use one of: "passphrase", "kms-gcp"\n`, wrapType)
			os.Exit(1)
		}

		errorCheck("sealing vault", vault.Seal(boxer))
		errorCheck("writing wallet file", vault.WriteToFile(walletFile))

		vaultWrittenReport(walletFile, newKeys, len(vault.KeyBag.Keys))
	},
}

func init() {
	vaultCmd.AddCommand(vaultCreateCmd)

	vaultCreateCmd.Flags().IntP("keys", "k", 0, "Number of keypairs to create")
	vaultCreateCmd.Flags().BoolP("import", "i", false, "Whether to import keys instead of creating them. This takes precedence over --keys, and private keys will be inputted on the command line.")
	vaultCreateCmd.Flags().StringP("comment", "c", "", "Comment field in the vault's json file.")
	vaultCreateCmd.Flags().StringP("vault-type", "t", "passphrase", "Vault type. One of: passphrase, kms-gcp")

	for _, flag := range []string{"keys", "comment", "vault-type", "import"} {
		if err := viper.BindPFlag("vault-create-cmd-"+flag, vaultCreateCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}
