// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"github.com/eoscanada/eos-go/msig"
	"github.com/spf13/cobra"
)

// msigUnapproveCmd represents the `eosio.msig::unapprove` command
var msigUnapproveCmd = &cobra.Command{
	Use:   "unapprove [proposer] [proposal name] [actor@permission]",
	Short: "Unapprove a transaction in the eosio.msig contract",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		proposer := toAccount(args[0], "proposer")
		proposalName := toName(args[1], "proposal name")
		requested, err := permissionToPermissionLevel(args[2])
		errorCheck("requested permission", err)

		pushEOSCActions(api,
			msig.NewUnapprove(proposer, proposalName, requested),
		)
	},
}

func init() {
	msigCmd.AddCommand(msigUnapproveCmd)
}
