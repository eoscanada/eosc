// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/spf13/cobra"
)

// txCmd represents the tx command
var txCmd = &cobra.Command{
	Use:   "tx",
	Short: "Transactions-related commands, like signing, pushing, reading, etc..",
}

func init() {
	RootCmd.AddCommand(txCmd)
}
