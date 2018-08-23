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
	Use:   "post [from_account] [message]",
	Short: "Post a message",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		targetAccount := toAccount(viper.GetString("forum-cmd-target-contract"), "--target-contract")
		accountName := toAccount(args[0], "from account")
		message := args[1]
		title := viper.GetString("forum-post-cmd-title")
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

		action := forum.NewPost(accountName, newUUID, title, message, replyTo, replyToUUID, certify, metadata)
		action.Account = targetAccount

		api := getAPI()
		pushEOSCActions(api, action)
	},
}

func init() {
	forumCmd.AddCommand(forumPostCmd)

	forumPostCmd.Flags().StringP("title", "", "", "The title for the post. None by default")
	forumPostCmd.Flags().BoolP("certify", "", false, "Certify that the contents of this message is true. See corresponding Ricardian Contract.")
	forumPostCmd.Flags().StringP("type", "", "chat", "Message type (added to json_metadata)")
	forumPostCmd.Flags().StringP("metadata", "", "", "Additional metadata. Must be JSON-encoded. If present, takes precedences over --type")
	forumPostCmd.Flags().StringP("reply-to", "", "", "Account name to reply to")
	forumPostCmd.Flags().StringP("reply-to-uuid", "", "", "UUID from a previous post from the --repy-to account.")

	for _, flag := range []string{"title", "certify", "type", "metadata", "reply-to", "reply-to-uuid"} {
		if err := viper.BindPFlag("forum-post-cmd-"+flag, forumPostCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}

}
