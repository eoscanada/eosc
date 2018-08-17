// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// forumCmd represents the forum command
var forumCmd = &cobra.Command{
	Use:   "forum",
	Short: "EOS Forum And Referendum interactions",
}

func init() {
	RootCmd.AddCommand(forumCmd)

	forumCmd.PersistentFlags().StringP("target-contract", "", "eosforumdapp", "Target account hosting the eosio.forum code")

	for _, flag := range []string{"target-contract"} {
		if err := viper.BindPFlag("forum-cmd-"+flag, forumCmd.PersistentFlags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}
