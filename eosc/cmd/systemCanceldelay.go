// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var systemCanceldelayCmd = &cobra.Command{
	Use:   `canceldelay`,
	Short: "Cancel a deferred transaction. Use 'eosc tx cancel' instead",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use `eosc tx cancel` instead.")
	},
}

func init() {
	systemCmd.AddCommand(systemCanceldelayCmd)
}
