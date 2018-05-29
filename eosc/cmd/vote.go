package cmd

import (
	"fmt"
	"os"

	"sort"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var voteCmd = &cobra.Command{
	Use:   "vote [voter name] [producer list]",
	Short: "Command to vote for block producers",
	Long:  `Command to vote for block producers`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var producerNames = make([]eos.AccountName, 0, 0)

		producerStringNames := args[1:]
		sort.Strings(producerStringNames)
		for _, producerString := range producerStringNames {
			producerNames = append(producerNames, eos.AccountName(producerString))
		}

		if len(producerNames) == 0 {
			fmt.Printf("No producer provided")
			os.Exit(1)
		}

		api, err := api()
		if err != nil {
			fmt.Printf("Error initiating api, %s\n", err.Error())
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
}
