// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/forum"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var forumUnPostCmd = &cobra.Command{
	Use:   "unpost [poster] [post_uuid]",
	Short: "Removes a given post.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		targetAccount := toAccount(viper.GetString("forum-cmd-target-contract"), "--target-contract")

		poster := toAccount(args[0], "from poster")
		postUUID := args[1]

		action := forum.NewUnPost(poster, postUUID)
		action.Account = targetAccount

		api := getAPI()
		pushEOSCActions(api, action)
	},
}

func init() {
	forumCmd.AddCommand(forumUnPostCmd)
}
