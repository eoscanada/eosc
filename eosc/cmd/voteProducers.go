package cmd

import (
	"context"
	"fmt"
	"sort"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var voteProducersCmd = &cobra.Command{
	Use:   "producers [voter name] [producer list]",
	Short: "Cast your vote for 1 to 30 producers. View them with 'list-producers'.",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		voterName := toAccount(args[0], "voter name")

		producerStringNames := args[1:]
		sort.Strings(producerStringNames)

		var producerNames []eos.AccountName
		for _, producerString := range producerStringNames {
			producerNames = append(producerNames, toAccount(producerString, "producer list"))
		}

		api := getAPI()

		fmt.Printf("Voter [%s] voting for: %s\n", voterName, producerNames)
		pushEOSCActions(context.Background(), api,
			system.NewVoteProducer(
				voterName,
				"",
				producerNames...,
			),
		)
	},
}

func init() {
	voteCmd.AddCommand(voteProducersCmd)
}
