package cmd

import (
	"fmt"

	eos "github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
)

var voteStatusCmd = &cobra.Command{
	Use:   "status [voter name]",
	Short: "Display the current vote status for a given account.",
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
		for _, info := range voterInfo {
			if info.Owner == voterName {
				found = true

				fmt.Println("Voter: ", info.Owner)
				fmt.Println("Is Proxy: ", info.IsProxy)
				if info.Proxy != "" {
					fmt.Println("Voting via proxy: ", info.Proxy)
					fmt.Println("Proxied vote weight: ", info.ProxiedVoteWeight)
				} else {
					fmt.Println("Producers list: ", info.Producers)
					fmt.Println("Staked amount: ", info.Staked)
					fmt.Println("Last vote weight: ", info.LastVoteWeight)
				}
			}
		}
		if !found {
			errorCheck("voter_info", fmt.Errorf("not found"))
		}
	},
}

func init() {
	voteCmd.AddCommand(voteStatusCmd)
}
