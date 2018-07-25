// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/spf13/cobra"
)

// msigCmd represents the msig command
var msigCmd = &cobra.Command{
	Use:   "multisig",
	Short: "eosio.msig contract interactions",
}

func init() {
	RootCmd.AddCommand(msigCmd)
}
