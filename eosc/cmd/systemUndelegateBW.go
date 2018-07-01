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
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var systemUndelegateBWCmd = &cobra.Command{
	Use:   "undelegatebw [from] [receiver] [network bw unstake qty] [cpu bw unstake qty]",
	Short: "Undelegate some CPU and Network bandwidth.",
	Long: `Undelegate some CPU and Network bandwidth.

When undelegating bandwidth, a "refund" action will automatically be
triggered and delayed for 72 hours.  This means it takes 3 days for
you to get your EOS back and being able to transfer it. However, your
voting power is immediately altered.

See also: the "system delegatebw" command.
`,
	Args: cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		from := toAccount(args[0], "from")
		receiver := toAccount(args[1], "receiver")
		netStake, err := eos.NewEOSAssetFromString(args[2])
		errorCheck(`"network bw unstake qty" invalid`, err)
		cpuStake, err := eos.NewEOSAssetFromString(args[3])
		errorCheck(`"cpu bw unstake qty" invalid`, err)

		api := apiWithWallet()

		pushEOSCActions(api, system.NewUndelegateBW(from, receiver, cpuStake, netStake))
	},
}

func init() {
	systemCmd.AddCommand(systemUndelegateBWCmd)
}
