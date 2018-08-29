// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"fmt"
	"os"

	"github.com/eoscanada/eos-go/msig"
	"github.com/spf13/cobra"
)

// msigApproveCmd represents the `eosio.msig::approve` command
var msigApproveCmd = &cobra.Command{
	Use:   "approve [proposer] [proposal name] [approver[@active]]",
	Short: "Approve a transaction in the eosio.msig contract",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		proposer := toAccount(args[0], "proposer")
		proposalName := toName(args[1], "proposal name")
		requested, err := permissionToPermissionLevel(args[2])
		if err != nil {
			fmt.Printf("Error with requested permission: %s\n", err)
			os.Exit(1)
		}

		pushEOSCActions(api,
			msig.NewApprove(proposer, proposalName, requested),
		)
	},
}

func init() {
	msigCmd.AddCommand(msigApproveCmd)
}
