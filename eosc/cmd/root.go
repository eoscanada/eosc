package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Version represents the eosc command version
var Version string

// RootCmd represents the eosc command
var RootCmd = &cobra.Command{
	Use:   "eosc",
	Short: "eosc is an EOS command-line Swiss Army knife",
	Long: `eosc is a command-line Swiss Army knife for EOS - by EOS Canada

It contains a Vault (or a wallet), a tool for voting, tools for end
users and tools for Block Producers.

The 'vault' acts as a keosd-compatible wallet (the one developed by
Block.one), while allowing you to manage your keys, and unlock it from
the command line.

Source code is available at: https://github.com/eoscanada/eosc
`,
}

// Execute executes the configured RootCmd
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
	RootCmd.PersistentFlags().StringP("api-url", "u", "https://mainnet.eoscanada.com", "API endpoint of eos.io blockchain node")
	RootCmd.PersistentFlags().StringSliceP("permission", "p", []string{}, "Permission to sign transactions with. Optionally specify more than one, or separate by comma")
	RootCmd.PersistentFlags().StringP("kms-gcp-keypath", "", "", "Path to the cryptoKeys within a keyRing on GCP")
	RootCmd.PersistentFlags().StringP("write-transaction", "", "", "Do not broadcast the transaction produced, but write it in json to the given filename instead.")
	RootCmd.PersistentFlags().StringP("offline-head-block", "", "", "Provide a recent block ID (long-form hex) for TaPoS. Use all --offline options to sign transactions offline.")
	RootCmd.PersistentFlags().StringP("offline-chain-id", "", "", "Chain ID to sign transaction with. Use all --offline- options to sign transactions offline.")
	RootCmd.PersistentFlags().StringSliceP("offline-sign-key", "", []string{}, "Public key to use to sign transaction. Must be in your vault or wallet. Use all --offline- options to sign transactions offline.")
	RootCmd.PersistentFlags().BoolP("skip-sign", "", false, "Do not sign the transaction. Use with --write-transaction.")
	RootCmd.PersistentFlags().IntP("expiration", "", 30, "Set time before transaction expires, in seconds. Defaults to 30 seconds.")
	RootCmd.PersistentFlags().IntP("delay-sec", "", 0, "Set time to wait before transaction is executed, in seconds. Defaults to 0 second.")
	RootCmd.PersistentFlags().BoolP("sudo-wrap", "", false, "Wrap the transaction in a eosio.sudo exec. Useful to BPs, with --write-transaction and --skip-sign to then submit as a multisig proposition.")

	for _, flag := range []string{"vault-file", "api-url", "kms-gcp-keypath", "wallet-url", "permission", "expiration", "delay-sec", "write-transaction", "skip-sign", "offline-head-block", "offline-chain-id", "offline-sign-key", "sudo-wrap"} {
		if err := viper.BindPFlag("global-"+flag, RootCmd.PersistentFlags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}

func initConfig() {
	viper.SetEnvPrefix("EOSC")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
}
