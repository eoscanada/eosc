package cmd

import (
	"fmt"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var voteRecastCmd = &cobra.Command{
	Use:   "recast [voter name]",
	Short: "Recast your vote for the same producers",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()
		voterName := toAccount(args[0], "voter name")

		response, err := api.GetTableRows(
			eos.GetTableRowsRequest{
				Code:       "eosio",
				Scope:      "eosio",
				Table:      "voters",
				JSON:       true,
				LowerBound: string(voterName),
				Limit:      1,
			},
		)
		errorCheck("get table row", err)

		var voterInfo []eos.VoterInfo
		err = response.JSONToStructs(&voterInfo)
		errorCheck("reading voter_info", err)

		var found bool
		var producerNames []eos.AccountName
		for _, info := range voterInfo {
			if info.Owner == voterName {
				found = true
				producerNames = info.Producers
			}
		}
		if !found {
			errorCheck("voter_info", fmt.Errorf("not found"))
			return
		}

		fmt.Printf("Voter [%s] recasting vote for: %s\n", voterName, producerNames)
		pushEOSCActions(api,
			system.NewVoteProducer(
				voterName,
				"",
				producerNames...,
			),
		)
	},
}

func init() {
	voteCmd.AddCommand(voteRecastCmd)
}
