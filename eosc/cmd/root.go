package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Version string

var RootCmd = &cobra.Command{
	Use:   "eosc",
	Short: "eosc is an EOS command-line swiss knife",
	Long: `eosc is an EOS command-line swiss knife

It contains a Vault (or a wallet), a tool for voting, tools for end
users and tools for Block Producers.

It is developed by EOS Canada, a (candidate) Block Producer for the EOS
network. Source code is available at: https://github.com/eoscanada/eosc

The 'vault' acts as a keosd-compatible wallet (the one developed by
Block.one), while allowing you to manage your keys, and unlock it from
the command line.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringP("vault-file", "", "./eosc-vault.json", "Wallet file that contains encrypted key material")
	RootCmd.PersistentFlags().StringSliceP("wallet-url", "", []string{}, "Base URL to wallet endpoint. You can pass this multiple times to use the multi-signer (will use each wallet to sign multi-sig transactions).")
	RootCmd.PersistentFlags().StringP("api-url", "u", "https://mainnet.eoscanada.com", "api endpoint of eos.io block chain node ")
	RootCmd.PersistentFlags().StringSliceP("permission", "p", []string{}, "Permission to sign transactions with. Optionally specify more than one, or separate by comma")
	RootCmd.PersistentFlags().StringP("kms-gcp-keypath", "", "", "Path to the cryptoKeys within a keyRing on GCP")

	for _, flag := range []string{"vault-file", "api-url", "kms-gcp-keypath", "wallet-url", "permission"} {
		if err := viper.BindPFlag("global-"+flag, RootCmd.PersistentFlags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}

func initConfig() {
	viper.SetEnvPrefix("EOSC")
	viper.AutomaticEnv()
}
