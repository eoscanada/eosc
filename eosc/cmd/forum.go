// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/spf13/cobra"
)

// forumCmd represents the forum command
var forumCmd = &cobra.Command{
	Use:   "forum",
	Short: "Forum messaging interactions",
}

func init() {
	RootCmd.AddCommand(forumCmd)
}
