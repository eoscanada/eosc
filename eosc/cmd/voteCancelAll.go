package cmd

import (
	"fmt"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var voteCancelAllCmd = &cobra.Command{
	Use:   "cancel-all [voter name]",
	Short: "Cancel all votes currently cast for producers/delegated to a proxy.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		voterName := toAccount(args[0], "voter name")

		noProxy := eos.AccountName("")
		var noVotes []eos.AccountName
		pushEOSCActions(api,
			system.NewVoteProducer(
				voterName,
				noProxy,
				noVotes...,
			),
		)

		fmt.Printf("Consider using `eosc vote status %s` to confirm it has been applied.\n", voterName)
	},
}

func init() {
	voteCmd.AddCommand(voteCancelAllCmd)
}
