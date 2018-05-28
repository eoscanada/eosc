package cmd

import (
	"fmt"

	"github.com/apex/log"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var voteCmd = &cobra.Command{
	Use:   "vote [voter name] [producer list]",
	Short: "Command to vote for block producers",
	Long:  `Command to vote for block producers`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		api := eos.New(
			viper.GetString("api-address"),
		)

		api.Debug = true

		keyBag := eos.NewKeyBag()
		err := keyBag.ImportFromFile(viper.GetString("keys-file"))
		if err != nil {
			log.Errorf("Error loading keys, %s\n", err.Error())
			return
		}
		fmt.Println("Keys file loaded")
		api.SetSigner(keyBag)

		var producerNames = make([]eos.AccountName, 0, 0)

		for _, producerString := range args[1:] {
			producerNames = append(producerNames, eos.AccountName(producerString))
		}

		if len(producerNames) == 0 {
			log.Error("No producer provided")
			return
		}

		voterName := args[0]
		fmt.Printf("Voter [%s] voting for: %s\n", voterName, producerNames)
		out, err := api.SignPushActions(
			system.NewVoteProducer(
				eos.AccountName(voterName),
				"",
				producerNames...,
			),
		)

		if err != nil {
			log.Errorf("Error during vote, %s\n", err.Error())
			return
		}

		fmt.Println("Vote sent to chain. ", out)
	},
}

func init() {

	RootCmd.AddCommand(voteCmd)

	voteCmd.Flags().StringP("api-address", "", "", "file containing signing key")
	voteCmd.Flags().StringP("keys-file", "", "", "file containing signing key")
	voteCmd.MarkFlagRequired("api-address")
	voteCmd.MarkFlagRequired("keys-file")

	if err := viper.BindPFlag("api-address", voteCmd.Flags().Lookup("api-address")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("keys-file", voteCmd.Flags().Lookup("keys-file")); err != nil {
		panic(err)
	}

}
