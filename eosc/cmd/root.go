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
the command line.

Version ` + Version,
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
	RootCmd.PersistentFlags().StringP("api-url", "u", "http://localhost:8888", "server api address")
	RootCmd.PersistentFlags().StringP("kms-gcp-keypath", "", "", "Path to the cryptoKeys within a keyRing on GCP")

	for _, flag := range []string{"vault-file", "api-url", "kms-gcp-keypath"} {
		if err := viper.BindPFlag("global-"+flag, RootCmd.PersistentFlags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}

func initConfig() {
	viper.SetEnvPrefix("EOSC")
	viper.AutomaticEnv()
}
