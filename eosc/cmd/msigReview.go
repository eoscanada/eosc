// Copyright © 2018 EOS Canada <alex@eoscanada.com>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	eos "github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
)

var msigReviewCmd = &cobra.Command{
	Use:   "review [proposer] [proposal name]",
	Short: "Review a proposal in the eosio.msig contract",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		proposer := eos.AccountName(args[0])
		proposalName := eos.Name(args[1])

		response, err := api.GetTableRows(
			eos.GetTableRowsRequest{
				Code:       "eosio.msig",
				Scope:      string(proposer),
				Table:      "proposal",
				JSON:       true,
				LowerBound: string(proposalName),
				Limit:      1,
			},
		)
		errorCheck("get table row", err)

		var transactions []struct {
			ProposalName eos.Name     `json:"proposal_name"`
			Transaction  eos.HexBytes `json:"packed_transaction"`
		}
		err = response.JSONToStructs(&transactions)
		errorCheck("reading proposed transactions", err)

		var transaction *eos.Transaction
		for _, tx := range transactions {
			if tx.ProposalName == proposalName {
				err := eos.UnmarshalBinary(tx.Transaction, &transaction)
				errorCheck("unmarshal binary data", err)
			}
		}
		if transaction == nil {
			errorCheck("transaction:", fmt.Errorf("not found"))
		}

		fmt.Println("Look for:", proposalName)

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(response)
	},
}

func init() {
	msigCmd.AddCommand(msigReviewCmd)
}
