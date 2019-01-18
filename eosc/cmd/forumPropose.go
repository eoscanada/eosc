// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"errors"
	"time"

	"github.com/eoscanada/eos-go"

	"github.com/eoscanada/eos-go/forum"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/sjson"
)

var forumProposeCmd = &cobra.Command{
	Use:   "propose [proposer] [proposal_name] [title] [proposal_expiration_date]",
	Short: "Submit a proposition for votes",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		targetAccount := toAccount(viper.GetString("forum-cmd-target-contract"), "--target-contract")

		proposer := toAccount(args[0], "proposer")
		proposalName := toName(args[1], "proposal_name")
		title := args[2]

		expiresAtStr := args[3]
		var expiresAt eos.JSONTime
		var err error
		if expiresAtStr != "" {
			expiresAt, err = eos.ParseJSONTime(expiresAtStr)
			errorCheck("no valid proposal expiration date provided. Must be set as ISO-8601 format. Maximum 6 months in the future.", err)
		}
		if expiresAt.Before(time.Now()) {
			errorCheck("proposal expiration date must in the future", errors.New("provided time is in the past"))
		}

		proposalJSON := viper.GetString("forum-propose-cmd-json")
		content := viper.GetString("forum-propose-cmd-content")
		jsonType := viper.GetString("forum-propose-cmd-type")
		if proposalJSON == "" && content != "" {
			proposalJSON = "{}"
		}
		proposalJSON, err = sjson.Set(proposalJSON, "content", content)
		// Defaults JSON schema type to `bps-proposal-v1`
		proposalJSON, err = sjson.Set(proposalJSON, "type", jsonType)
		errorCheck("setting content in json", err)

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
	forumProposeCmd.Flags().String("type", "bps-proposal-v1", "The JSON schema of the proposal, set as a `type` in the JSON payload.")
}
