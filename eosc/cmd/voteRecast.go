package cmd

import (
	"fmt"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var voteRecastCmd = &cobra.Command{
	Use:   "recast [voter name]",
	Short: "Recast your vote for the same producers or proxy.",
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

		var voterInfos []eos.VoterInfo
		err = response.JSONToStructs(&voterInfos)
		errorCheck("reading voter_info", err)

		found := false
		for _, info := range voterInfos {
			if info.Owner == voterName {
				found = true
				if info.Proxy != "" {
					fmt.Printf("Voter [%s] recasting vote via proxy: %s\n", voterName, info.Proxy)
				} else {
					voterPrefix := ""
					if info.IsProxy != 0 {
						voterPrefix = "Proxy "
					}
					producersList := "no producer"
					if len(info.Producers) >= 1 {
						producersList = fmt.Sprint(info.Producers)
					}
					fmt.Printf("%sVoter [%s] recasting vote for: %s\n", voterPrefix, voterName, producersList)
				}
				pushEOSCActions(api,
					system.NewVoteProducer(
						voterName,
						info.Proxy,
						info.Producers...,
					),
				)
			}
		}
		if !found {
			errorCheck("vote recast", fmt.Errorf("unable to recast vote as no existing vote was found"))
		}
	},
}

func init() {
	voteCmd.AddCommand(voteRecastCmd)
}
