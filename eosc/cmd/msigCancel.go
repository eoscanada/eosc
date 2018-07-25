// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/msig"
	"github.com/spf13/cobra"
)

// msigCancelCmd represents the `eosio.msig::cancel` command
var msigCancelCmd = &cobra.Command{
	Use:   "cancel [proposer] [proposal name] [canceler]",
	Short: "Cancel a transaction in the eosio.msig contract",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		proposer := toAccount(args[0], "proposer")
		proposalName := toName(args[1], "proposal name")
		canceler := toAccount(args[2], "canceler")

		pushEOSCActions(api,
			msig.NewCancel(proposer, proposalName, canceler),
		)
	},
}

func init() {
	msigCmd.AddCommand(msigCancelCmd)
}
