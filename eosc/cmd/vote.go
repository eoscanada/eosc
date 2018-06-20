package cmd

import (
	"github.com/spf13/cobra"
)

var voteCmd = &cobra.Command{
	Use:   "vote",
	Short: "Command to vote for block producers or proxy",
}

func init() {
	RootCmd.AddCommand(voteCmd)
}
