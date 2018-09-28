// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/spf13/cobra"
)

var devScaffoldCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "Scaffolds an EOS Smart Contract project.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	devCmd.AddCommand(devScaffoldCmd)
}
