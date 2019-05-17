// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/spf13/cobra"
)

var rexCmd = &cobra.Command{
	Use:   "rex",
	Short: "EOS REX interactions",
}

func init() {
	RootCmd.AddCommand(rexCmd)
}
