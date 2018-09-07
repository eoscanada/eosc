// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"fmt"

	eos "github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
)

var msigStatusCmd = &cobra.Command{
	Use:   "status [proposer] [proposal name]",
	Short: "Shows the status of a proposal and its approvals in the eosio.msig contract.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		proposer := toAccount(args[0], "proposer")
		proposalName := toName(args[1], "proposal name")

		response, err := api.GetTableRows(
			eos.GetTableRowsRequest{
				Code:       "eosio.msig",
				Scope:      string(proposer),
				Table:      "approvals",
				JSON:       true,
				LowerBound: string(proposalName),
				// UpperBound: string(proposalName),
				Limit: 1,
			},
		)
		errorCheck("get table row", err)

		var approvalsInfo []struct {
			ProposalName       eos.Name              `json:"proposal_name"`
			RequestedApprovals []eos.PermissionLevel `json:"requested_approvals"`
			ProvidedApprovals  []eos.PermissionLevel `json:"provided_approvals"`
		}
		err = response.JSONToStructs(&approvalsInfo)
		errorCheck("reading approvals_info", err)

		if len(approvalsInfo) == 1 {
			fmt.Println("Proposal name:", approvalsInfo[0].ProposalName)
			fmt.Println("Requested approvals:", approvalsInfo[0].RequestedApprovals)
			fmt.Println("Provided approvals:", approvalsInfo[0].ProvidedApprovals)
		} else {
			fmt.Printf("Proposal %s not found.", string(proposalName))
		}
	},
}

func init() {
	msigCmd.AddCommand(msigStatusCmd)
}
