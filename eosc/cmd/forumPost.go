// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"encoding/json"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/forum"
	"github.com/pborman/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var forumPostCmd = &cobra.Command{
	Use:   "post [poster] [content]",
	Short: "Post a message",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		targetAccount := toAccount(viper.GetString("forum-cmd-target-contract"), "--target-contract")

		poster := toAccount(args[0], "from poster")
		content := args[1]

		certify := viper.GetBool("forum-post-cmd-certify")
		newUUID := uuid.New()

		metadata := viper.GetString("forum-post-cmd-metadata")
		if metadata != "" {
			var dump interface{}
			err := json.Unmarshal([]byte(metadata), &dump)
			errorCheck("--metadata is not valid JSON", err)
		} else {
			metadataBytes, _ := json.Marshal(map[string]interface{}{
				"type": viper.GetString("forum-post-cmd-type"),
			})
			metadata = string(metadataBytes)
		}

		replyTo := eos.AccountName(viper.GetString("forum-post-cmd-reply-to"))
		if len(replyTo) != 0 {
			_ = toAccount(string(replyTo), "--reply-to") // only check for errors
		}

		replyToUUID := viper.GetString("forum-post-cmd-reply-to-uuid")

		action := forum.NewPost(poster, newUUID, content, replyTo, replyToUUID, certify, metadata)
		action.Account = targetAccount

		api := getAPI()
		pushEOSCActions(api, action)
	},
}

func init() {
	forumCmd.AddCommand(forumPostCmd)

	forumPostCmd.Flags().Bool("certify", false, "Certify that the contents of this message is true. See corresponding Ricardian Contract.")
	forumPostCmd.Flags().String("type", "chat", "Message type (added to json_metadata)")
	forumPostCmd.Flags().String("metadata", "", "Additional metadata. Must be JSON-encoded. If present, takes precedences over --type")
	forumPostCmd.Flags().String("reply-to", "", "Account name to reply to")
	forumPostCmd.Flags().String("reply-to-uuid", "", "UUID from a previous post from the --repy-to account.")
}
