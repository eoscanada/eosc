// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/spf13/cobra"
)

var devInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes the EOS Smart Contract developer environment.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		sanityCheckDevEnvironment()
	},
}

func init() {
	devCmd.AddCommand(devInitCmd)
}
