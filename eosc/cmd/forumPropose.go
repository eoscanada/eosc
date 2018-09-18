// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go"

	"github.com/eoscanada/eos-go/forum"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/sjson"
)

var forumProposeCmd = &cobra.Command{
	Use:   "propose [proposer] [proposal_name] [title]",
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

		expiresAt, err := eos.ParseJSONTime(viper.GetString("forum-propose-cmd-proposal-expires-at"))
		errorCheck("unable to parse expiration time", err)

		action := forum.NewPropose(proposer, proposalName, title, proposalJSON, expiresAt)
		action.Account = targetAccount

		api := getAPI()
		pushEOSCActions(api, action)
	},
}

func init() {
	forumCmd.AddCommand(forumProposeCmd)

	forumProposeCmd.Flags().String("content", "", "Markdown 'content' to be injected in the JSON (whether you propose a --json or not).")
	forumProposeCmd.Flags().String("json", "", "Proposal JSON body.")
	forumProposeCmd.Flags().String("proposal-expires-at", "", "Time at which the proposal expires (maximum 6 months in the future). This must be formatted as ISO-8601 datetime.")

	for _, flag := range []string{"content", "json", "proposal-expires-at"} {
		if err := viper.BindPFlag("forum-propose-cmd-"+flag, forumProposeCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}
