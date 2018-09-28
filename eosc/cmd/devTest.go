// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/spf13/cobra"
)

var devTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Runs the test suite of an EOS Smart Contract project.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	devCmd.AddCommand(devTestCmd)
}
