// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/forum"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var forumVoteCmd = &cobra.Command{
	Use:   "vote [voter] [proposition] [vote_value]",
	Short: "Submit a vote from [voter] on the [proposition] with a [vote_value] agreed in the proposition.",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {

		accountName := toAccount(args[0], "voter")
		proposition := args[1]
		vote := args[2]
		propositionHash := viper.GetString("forum-vote-cmd-hash")

		api := getAPI()
		pushEOSCActions(api,
			forum.NewVote(accountName, proposition, propositionHash, vote),
		)
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
