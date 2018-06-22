package cmd

import (
	"fmt"
	"os"
	"sort"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var voteProducersCmd = &cobra.Command{
	Use:   "producers [voter name] [producer list]",
	Short: "Cast your vote for 1 to 30 producers. View them with 'list'",
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

		api := apiWithWallet()

		voterName := args[0]

		fmt.Printf("Voter [%s] voting for: %s\n", voterName, producerNames)
		pushEOSCActions(api,
			system.NewVoteProducer(
				eos.AccountName(voterName),
				"",
				producerNames...,
			),
		)
	},
}

func init() {
	voteCmd.AddCommand(voteProducersCmd)
}
