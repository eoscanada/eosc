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
	Use:   "vote [voter] [proposal_name] [vote_value]",
	Short: "Submit a vote from [voter] on [proposal_name] with a [vote_value].",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		targetAccount := toAccount(viper.GetString("forum-cmd-target-contract"), "--target-contract")

		voter := toAccount(args[0], "voter")
		proposalName := toName(args[1], "proposal_name")

		// TODO: in a func
		vote := args[2]
		if vote == "yes" {
			vote = "1"
		}
		if vote == "no" {
			vote = "0"
		}
		voteValue, err := strconv.ParseUint(vote, 10, 8)
		errorCheck("expected an integer for vote_value", err)
		if voteValue > 255 {
			errorCheck("vote value cannot exceed 255", fmt.Errorf("vote value too high: %d", voteValue))
		}

		json := viper.GetString("forum-vote-cmd-json")

		action := forum.NewVote(voter, proposalName, uint8(voteValue), json)
		action.Account = targetAccount

		api := getAPI()
		pushEOSCActions(api, action)
	},
}

func init() {
	forumCmd.AddCommand(forumVoteCmd)

	forumVoteCmd.Flags().String("json", "", "Optional JSON attached to the vote.")
}
