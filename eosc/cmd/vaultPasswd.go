// Copyright © 2020 dfuse Platform inc <info@dfuse.io>

package cmd

import (
	"fmt"

	"github.com/eoscanada/eosc/cli"
	"github.com/eoscanada/eosc/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// vaultPasswdCmd represents the serve command
var vaultPasswdCmd = &cobra.Command{
	Use:   "serve",
	Short: "Change the passphrase on the vault",
	Run: func(cmd *cobra.Command, args []string) {
		v := mustGetWallet()

		if v.SecretBoxWrap != "passphrase" {
			errorCheck("only passphrase vaults supported to change passphrase", fmt.Errorf("not supported"))
		}

		password, err := cli.GetEncryptPassphrase()
		errorCheck("password input", err)

		boxer := vault.NewPassphraseBoxer(password)

		v.PrintPublicKeys()

		errorCheck("sealing vault", v.Seal(boxer))
		errorCheck("writing wallet file", v.WriteToFile(viper.GetString("global-vault-file")))
	},
}

func init() {
	vaultCmd.AddCommand(vaultPasswdCmd)
}
