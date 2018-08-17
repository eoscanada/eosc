// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/forum"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/sjson"
)

var forumProposeCmd = &cobra.Command{
	Use:   "propose [proposer] [proposal_name] [title string]",
	Short: "Submit a proposition for votes",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {

		targetAccount := toAccount(viper.GetString("forum-cmd-target-contract"), "--target-contract")
		proposer := toAccount(args[0], "proposer")
		proposalName := toName(args[1], "proposal_name")
		title := args[2]
		proposalJSON := viper.GetString("forum-propose-cmd-json")
		content := viper.GetString("forum-propose-cmd-content")
		if proposalJSON == "" && content != "" {
			proposalJSON = "{}"
		}
		proposalJSON, err := sjson.Set(proposalJSON, "content", content)
		errorCheck("setting content in json", err)

		action := forum.NewPropose(proposer, proposalName, title, proposalJSON)
		action.Account = targetAccount

		api := getAPI()
		pushEOSCActions(api, action)
	},
}

func init() {
	forumCmd.AddCommand(forumProposeCmd)

	forumProposeCmd.Flags().StringP("content", "", "", "Markdown 'content' to be injected in the JSON (whether you propose a --json or not)")
	forumProposeCmd.Flags().StringP("json", "", "", "JSON")

	for _, flag := range []string{"content", "json"} {
		if err := viper.BindPFlag("forum-propose-cmd-"+flag, forumProposeCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}
