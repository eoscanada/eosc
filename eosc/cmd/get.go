// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import "github.com/spf13/cobra"

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Fetch information from the blockchain",
}

func init() {
	RootCmd.AddCommand(getCmd)
}
