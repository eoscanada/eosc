// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/forum"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var forumUnVoteCmd = &cobra.Command{
	Use:   "unvote [voter] [proposal_name]",
	Short: "Cancels a vote for a given proposal.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		targetAccount := toAccount(viper.GetString("forum-cmd-target-contract"), "--target-contract")

		voter := toAccount(args[0], "voter")
		proposalName := toName(args[1], "proposal_name")

		action := forum.NewUnVote(voter, proposalName)
		action.Account = targetAccount

		api := getAPI()
		pushEOSCActions(api, action)
	},
}

func init() {
	forumCmd.AddCommand(forumUnVoteCmd)
}
