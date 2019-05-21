// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

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
				Table: "proposal",
				JSON:  true,
			},
		)
		errorCheck("get table row", err)

		var proposals []proposalRow
		err = response.JSONToStructs(&proposals)
		errorCheck("reading proposal list", err)

		if printJSON, _ := cmd.Flags().GetBool("json"); printJSON == true {
			data, err := json.MarshalIndent(proposals, "", "  ")
			errorCheck("json marshal", err)
			fmt.Println(string(data))
			return
		}

		if len(proposals) == 0 {
			errorCheck("No multisig proposal found", fmt.Errorf("not found"))
		} else {
			fmt.Println("All active proposals")
			fmt.Println("---------------------")
			for _, proposal := range proposals {
				fmt.Println("Proposal name:", proposal.ProposalName)
			}
			fmt.Println("---------------------")
		}

		response, err = api.GetTableRows(
			eos.GetTableRowsRequest{
				Code:  "eosio.msig",
				Scope: string(proposer),
				Table: "approvals",
				JSON:  true,
			},
		)
		errorCheck("get table row", err)

		var approvals []approvalRow
		err = response.JSONToStructs(&approvals)
		errorCheck("reading approval_info list", err)

		if printJSON, _ := cmd.Flags().GetBool("json"); printJSON == true {
			data, err := json.MarshalIndent(approvals, "", "  ")
			errorCheck("json marshal", err)
			fmt.Println(string(data))
			return
		}

		for _, info := range approvals {
			info.Show()
			fmt.Println()
		}
	},
}

type proposalRow struct {
	ProposalName eos.Name     `json:"proposal_name"`
	Transaction  eos.HexBytes `json:"packed_transaction"`
}
type approvalRow struct {
	ProposalName       eos.Name              `json:"proposal_name"`
	RequestedApprovals []eos.PermissionLevel `json:"requested_approvals"`
	ProvidedApprovals  []eos.PermissionLevel `json:"provided_approvals"`
}

func (a approvalRow) Show() {
	fmt.Println("Proposal name:", a.ProposalName)
	fmt.Println("Requested approvals:")
	fmt.Print(formatAuths(a.RequestedApprovals))
	fmt.Println("Provided approvals:")
	fmt.Print(formatAuths(a.ProvidedApprovals))
	fmt.Println("---------------------")
}

func formatAuths(perms []eos.PermissionLevel) string {
	var out []string
	for _, perm := range perms {
		out = append(out, fmt.Sprintf("- %s@%s\n", perm.Actor, perm.Permission))
	}
	return strings.Join(out, "")
}

func init() {
	msigCmd.AddCommand(msigListCmd)
	msigListCmd.Flags().BoolP("json", "", false, "Display as JSON - useful to tally approvals")
}
