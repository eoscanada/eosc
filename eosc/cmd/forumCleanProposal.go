// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"strconv"

	"github.com/eoscanada/eos-go/forum"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var forumCleanProposalCmd = &cobra.Command{
	Use:   "clean-proposal [cleaner_account] [proposal_name] [max_count]",
	Short: "Cleans an expired proposal",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		targetAccount := toAccount(viper.GetString("forum-cmd-target-contract"), "--target-contract")

		cleaner := toAccount(args[0], "cleaner_account")
		proposalName := toName(args[1], "proposal_name")
		maxCount, err := strconv.ParseUint(args[2], 10, 64)
		errorCheck("Unable to parse max_count argument", err)

		action := forum.NewCleanProposal(cleaner, proposalName, maxCount)
		action.Account = targetAccount

		api := getAPI()
		pushEOSCActions(api, action)
	},
}

func init() {
	forumCmd.AddCommand(forumCleanProposalCmd)
}
