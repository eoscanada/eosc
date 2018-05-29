package cmd

import (
	"fmt"
	"os"

	"github.com/dgiagio/getpass"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
	eosvault "github.com/eoscanada/eosc/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var voteCmd = &cobra.Command{
	Use:   "vote [voter name] [producer list]",
	Short: "Command to vote for block producers",
	Long:  `Command to vote for block producers`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		walletFile := viper.GetString("vault-file")
		if _, err := os.Stat(walletFile); err != nil {
			fmt.Printf("Wallet file %q missing, \n", walletFile)
			os.Exit(1)
		}

		vault, err := eosvault.NewVaultFromWalletFile(walletFile)
		if err != nil {
			fmt.Printf("Error loading vault, %s\n", err)
			os.Exit(1)
		}

		passphrase, err := getpass.GetPassword("Enter passphrase to unlock vault: ")
		if err != nil {
			fmt.Println("ERROR reading passphrase:", err)
			os.Exit(1)
		}

		switch vault.SecretBoxWrap {
		case "passphrase":
			err = vault.OpenWithPassphrase(passphrase)
			if err != nil {
				fmt.Println("ERROR reading passphrase:", err)
				os.Exit(1)
			}
		default:
			fmt.Printf("ERROR unsupported secretbox wrapping method: %q\n", vault.SecretBoxWrap)
			os.Exit(1)
		}

		vault.PrintPublicKeys()

		api := eos.New(
			viper.GetString("api-address"),
		)

		api.SetSigner(vault.KeyBag)

		var producerNames = make([]eos.AccountName, 0, 0)

		for _, producerString := range args[1:] {
			producerNames = append(producerNames, eos.AccountName(producerString))
		}

		if len(producerNames) == 0 {
			fmt.Printf("No producer provided")
			os.Exit(1)
		}

		voterName := args[0]
		fmt.Printf("Voter [%s] voting for: %s\n", voterName, producerNames)
		_, err = api.SignPushActions(
			system.NewVoteProducer(
				eos.AccountName(voterName),
				"",
				producerNames...,
			),
		)

		if err != nil {
			fmt.Printf("Error during vote, %s\n", err.Error())
			os.Exit(1)
		}

		fmt.Println("Vote sent to chain.")
	},
}

func init() {

	RootCmd.AddCommand(voteCmd)

	voteCmd.Flags().StringP("api-address", "", "", "file containing signing key")
	voteCmd.MarkFlagRequired("api-address")

	if err := viper.BindPFlag("api-address", voteCmd.Flags().Lookup("api-address")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("keys-file", voteCmd.Flags().Lookup("keys-file")); err != nil {
		panic(err)
	}

}
