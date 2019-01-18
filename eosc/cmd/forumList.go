// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"encoding/json"
	"fmt"

	eos "github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var forumListCmd = &cobra.Command{
	Use:   "list",
	Short: "List forum proposals.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()
		targetAccount := toAccount(viper.GetString("forum-cmd-target-contract"), "--target-contract")

		proposerStr := viper.GetString("forum-list-cmd-from-proposer")
		var proposer eos.AccountName

		var err error
		var resp *eos.GetTableRowsResp
		if proposerStr != "" {
			proposer = toAccount(proposerStr, "--from-proposer")
			resp, err = api.GetTableRows(eos.GetTableRowsRequest{
				Code:       string(targetAccount),
				Scope:      string(targetAccount),
				Table:      string("proposal"),
				Index:      "sec", // Secondary index `by_proposer`
				KeyType:    "name",
				LowerBound: string(proposer),
				Limit:      1000,
				JSON:       true,
			})
			if err != nil {
				errorCheck(fmt.Sprintf("unable to get list of proposals from proposer %q", proposer), err)
			}
		} else {
			resp, err = api.GetTableRows(eos.GetTableRowsRequest{
				Code:  string(targetAccount),
				Scope: string(targetAccount),
				Table: string("proposal"),
				Limit: 1000,
				JSON:  true,
			})
			if err != nil {
				errorCheck("unable to get list of proposals", err)
			}

		}

		var proposals []struct {
			ProposalName eos.Name        `json:"proposal_name"`
			Proposer     eos.AccountName `json:"proposer"`
			Title        string          `json:"title"`
			ProposalJSON string          `json:"proposal_json"`
			CreatedAt    eos.JSONTime    `json:"created_at"`
			ExpiresAt    eos.JSONTime    `json:"expires_at"`
		}
		err = resp.JSONToStructs(&proposals)
		errorCheck("reading proposal list", err)

		if printJSON, _ := cmd.Flags().GetBool("json"); printJSON == true {
			data, err := json.MarshalIndent(proposals, "", "  ")
			errorCheck("json marshal", err)
			fmt.Println(string(data))
			return
		}

		found := false
		for _, proposal := range proposals {
			if proposerStr == "" || proposal.Proposer == proposer {
				fmt.Println("Proposal name: ", proposal.ProposalName)
				fmt.Println("Proposer: ", proposal.Proposer)
				fmt.Println("Title: ", proposal.Title)
				fmt.Println("JSON: ", proposal.ProposalJSON)
				fmt.Println("Created at: ", proposal.CreatedAt)
				fmt.Println("Expires at: ", proposal.ExpiresAt)
				fmt.Println()

				found = true
			}
		}
		if !found {
			errorCheck("no proposal found", fmt.Errorf("empty list"))
		}
	},
}

func init() {
	forumCmd.AddCommand(forumListCmd)

	forumListCmd.Flags().String("from-proposer", "", "Filter proposals only from proposer.")
	forumListCmd.Flags().Bool("json", false, "Output list as JSON")
}
