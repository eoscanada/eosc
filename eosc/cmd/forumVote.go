// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"fmt"
	"strconv"

	"github.com/eoscanada/eos-go/forum"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var forumVoteCmd = &cobra.Command{
	Use:   "vote [voter] [proposer] [proposal_name] [vote_value]",
	Short: "Submit a vote from [voter] on [proposer]'s [proposal_name] with a [vote_value] agreed in the proposition.",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		targetAccount := toAccount(viper.GetString("forum-cmd-target-contract"), "--target-contract")

		voter := toAccount(args[0], "voter")
		proposer := toAccount(args[1], "proposer")
		proposalName := toName(args[2], "proposal_name")

		// TODO: in a func
		vote := args[3]
		if vote == "yes" {
			vote = "1"
		}
		if vote == "no" {
			vote = "0"
		}
		voteValue, err := strconv.ParseInt(vote, 10, 64)
		errorCheck("expected an integer for vote_value", err)
		if voteValue > 255 {
			errorCheck("vote value cannot exceed 255", fmt.Errorf("vote value too high: %d", voteValue))
		}
		proposalHash := viper.GetString("forum-vote-cmd-hash")

		action := forum.NewVote(voter, proposer, proposalName, proposalHash, uint8(voteValue), "")
		action.Account = targetAccount

		api := getAPI()
		pushEOSCActions(api, action)
	},
}

func init() {
	forumCmd.AddCommand(forumVoteCmd)

	forumVoteCmd.Flags().StringP("hash", "", "", "Hash of the proposition, as defined by the proposition itself")

	for _, flag := range []string{"hash"} {
		if err := viper.BindPFlag("forum-vote-cmd-"+flag, forumVoteCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}
