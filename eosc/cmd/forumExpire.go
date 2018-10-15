// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/forum"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var forumExpireCmd = &cobra.Command{
	Use:   "expire [proposer] [proposal_name]",
	Short: "Allows the [proposer] to expire a proposal before its set expiration time.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		targetAccount := toAccount(viper.GetString("forum-cmd-target-contract"), "--target-contract")

		proposer := toAccount(args[0], "proposer")
		proposalName := toName(args[1], "proposal_name")

		action := forum.NewExpire(proposer, proposalName)
		action.Account = targetAccount

		api := getAPI()
		pushEOSCActions(api, action)
	},
}

func init() {
	forumCmd.AddCommand(forumExpireCmd)
}
