// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/spf13/cobra"
)

// forumCmd represents the forum command
var forumCmd = &cobra.Command{
	Use:   "forum",
	Short: "EOS Forum And Referendum interactions",
}

func init() {
	RootCmd.AddCommand(forumCmd)

	forumCmd.PersistentFlags().String("target-contract", "eosio.forum", "Target account hosting the eosio.forum code")
}
