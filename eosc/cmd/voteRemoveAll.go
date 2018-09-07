package cmd

import (
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var voteRemoveAllCmd = &cobra.Command{
	Use:   "remove-all [voter name]",
	Short: "Remove all votes currently casted - effectively voting for no one.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		voterName := toAccount(args[0], "voter name")

		var noVotes []eos.AccountName
		pushEOSCActions(api,
			system.NewVoteProducer(
				voterName,
				"",
				noVotes...,
			),
		)
	},
}

func init() {
	voteCmd.AddCommand(voteRemoveAllCmd)
}
