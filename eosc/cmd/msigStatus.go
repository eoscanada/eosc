// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"encoding/json"
	"fmt"

	eos "github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
)

var msigStatusCmd = &cobra.Command{
	Use:   "status [proposer] [proposal name]",
	Short: "Shows the status of a given proposal and its approvals in the eosio.msig contract.",
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
				Limit:      1,
			},
		)
		errorCheck("get table row", err)

		var approvals []approvalRow
		err = response.JSONToStructs(&approvals)
		errorCheck("reading approvals_info", err)

		var found bool
		for _, info := range approvals {
			if info.ProposalName == proposalName {
				found = true

				if printJSON, _ := cmd.Flags().GetBool("json"); printJSON == true {
					data, err := json.MarshalIndent(info, "", "  ")
					errorCheck("json marshal", err)
					fmt.Println(string(data))
				} else {
					fmt.Println("Proposer:", proposer)
					info.Show()
					fmt.Println("")
					fmt.Println("Review with:")
					fmt.Println("")
					fmt.Println("    eosc multisig review", proposer, proposalName)
					fmt.Println("")
					fmt.Println("Approve with:")
					fmt.Println("")
					fmt.Println("    eosc multisig approve", proposer, proposalName, "[your_account]")
					fmt.Println("")
					fmt.Println("Execute with:")
					fmt.Println("")
					fmt.Println("    eosc multisig exec", proposer, proposalName, "[your_account]")
					fmt.Println("")
				}
			}
		}
		if !found {
			errorCheck("multisig proposal", fmt.Errorf("not found"))
		}
	},
}

func init() {
	msigCmd.AddCommand(msigStatusCmd)
	msigStatusCmd.Flags().BoolP("json", "", false, "Display as JSON - useful to tally approvals")
}
