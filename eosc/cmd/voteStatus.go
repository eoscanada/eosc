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

		var voterInfos []eos.VoterInfo
		err = response.JSONToStructs(&voterInfos)
		errorCheck("reading voter_info", err)

		found := false
		for _, info := range voterInfos {
			if info.Owner == voterName {
				found = true
				fmt.Println("Voter: ", info.Owner)

				if info.IsProxy != 0 {
					fmt.Println("Registered as a proxy voter: true")
					fmt.Println("Proxied vote weight: ", info.ProxiedVoteWeight)
				} else {
					fmt.Println("Registered as a proxy voter: false")
				}

				if info.Proxy != "" {
					fmt.Println("Voting via proxy: ", info.Proxy)
					fmt.Println("Last vote weight: ", info.LastVoteWeight)

				} else {
					fmt.Println("Producers list: ", info.Producers)
					fmt.Println("Staked amount: ", info.Staked)
					fmt.Printf("Last vote weight: %f\n", info.LastVoteWeight)
				}
			}
		}
		if !found {
			errorCheck("vote status", fmt.Errorf("unable to find vote status for %s", voterName))
		}
	},
}

func init() {
	voteCmd.AddCommand(voteStatusCmd)
}
