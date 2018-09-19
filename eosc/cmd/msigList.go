// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"encoding/json"
	"fmt"

	eos "github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
)

var msigListCmd = &cobra.Command{
	Use:   "list [proposer]",
	Short: "Shows the list of all active proposals for a given proposer in the eosio.msig contract.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		proposer := toAccount(args[0], "proposer")

		response, err := api.GetTableRows(
			eos.GetTableRowsRequest{
				Code:  "eosio.msig",
				Scope: string(proposer),
				Table: "approvals",
				JSON:  true,
			},
		)
		errorCheck("get table row", err)

		var approvalsInfo []struct {
			ProposalName       eos.Name              `json:"proposal_name"`
			RequestedApprovals []eos.PermissionLevel `json:"requested_approvals"`
			ProvidedApprovals  []eos.PermissionLevel `json:"provided_approvals"`
		}
		err = response.JSONToStructs(&approvalsInfo)
		errorCheck("reading approvals_info list", err)

		if printJSON, _ := cmd.Flags().GetBool("json"); printJSON == true {
			data, err := json.MarshalIndent(approvalsInfo, "", "  ")
			errorCheck("json marshal", err)
			fmt.Println(string(data))
			return
		}

		for _, info := range approvalsInfo {
			fmt.Println("Proposal name:", info.ProposalName)
			fmt.Println("Requested approvals:", info.RequestedApprovals)
			fmt.Println("Provided approvals:", info.ProvidedApprovals)
			fmt.Println()
		}
		if len(approvalsInfo) == 0 {
			errorCheck("No multisig proposal found", fmt.Errorf("not found"))
		}
	},
}

func init() {
	msigCmd.AddCommand(msigListCmd)
	msigListCmd.Flags().BoolP("json", "", false, "Display as JSON - useful to tally approvals")
}
