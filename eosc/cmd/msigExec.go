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
	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/msig"
	"github.com/spf13/cobra"
)

// msigExecCmd represents the `eosio.msig::exec` command
var msigExecCmd = &cobra.Command{
	Use:   "exec [proposer] [proposal name] [executer]",
	Short: "Execute a transaction in the eosio.msig contract",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		api := apiWithWallet()

		proposer := eos.AccountName(args[0])
		proposalName := eos.Name(args[1])
		executer := eos.AccountName(args[2])

		pushEOSCActions(api,
			msig.NewExec(proposer, proposalName, executer),
		)
	},
}

func init() {
	msigCmd.AddCommand(msigExecCmd)
}
