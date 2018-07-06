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
	"github.com/spf13/viper"
)

var systemDelegateBWCmd = &cobra.Command{
	Use:   "delegatebw [from] [receiver] [network bw stake qty] [cpu bw stake qty]",
	Short: "Delegate some CPU and Network bandwidth, to yourself or others.",
	Long: `Delegate some CPU and Network bandwidth, to yourself or others.

Bandwidth on EOS allows you to submit transactions on the network.

Delegating bandwidth (oftentimes called "staking") and locking it up
for 72 hours has two effects: increasing your voting power, and
increasing the bandwidth you're allocated to use the network.

CPU bandwidth means the time taken by Block Producers (in micro or
milliseconds) to process your transaction.

Network bandwidth means the number of bytes your transaction consumes
when propagating your transaction on the network, and finally putting
it in a block.

Those two sorts of bandwidth have burst capacity, and once used, will
both re-increase as time goes by.

The --transfer option makes it so the receiver will be able to unstake
what was delegated to them, and receive the corresponding EOS back. It
is effectively transfering the coins to them.

Example use:

    eosc system delegatebw myaccount youraccount "1.0000 EOS" "2.0000 EOS"

Alternatively, you can use the simplified:

    eosc system delegatebw myaccount youraccount 1.0 2.0
`,
	Args: cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		from := toAccount(args[0], "from")
		receiver := toAccount(args[1], "receiver")
		netStake, err := eos.NewEOSAssetFromString(args[2])
		errorCheck(`"network bw stake qty" invalid`, err)
		cpuStake, err := eos.NewEOSAssetFromString(args[3])
		errorCheck(`"cpu bw stake qty" invalid`, err)
		transfer := viper.GetBool("system-delegatebw-cmd-transfer")

		api := getAPI()

		pushEOSCActions(api, system.NewDelegateBW(from, receiver, cpuStake, netStake, transfer))
	},
}

func init() {
	systemCmd.AddCommand(systemDelegateBWCmd)

	systemDelegateBWCmd.Flags().BoolP("transfer", "", false, "Transfer voting power and right to unstake EOS to receiver")

	for _, flag := range []string{"transfer"} {
		if err := viper.BindPFlag("system-delegatebw-cmd-"+flag, systemDelegateBWCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}
