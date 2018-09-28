// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/spf13/cobra"
)

var devBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds a Smart Contract",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	devCmd.AddCommand(devBuildCmd)
}
