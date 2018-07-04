// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

		api := apiWithWallet()
		pushEOSCActions(api,
			forum.NewPost(accountName, newUUID, title, message, replyTo, replyToUUID, certify, metadata),
		)
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
