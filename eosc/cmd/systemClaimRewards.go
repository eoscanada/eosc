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
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var systemClaimRewardsCmd = &cobra.Command{
	Use:   "claimrewards [owner]",
	Short: "Claim block production rewards. Once per day, don't forget it!",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		owner := toAccount(args[0], "owner")

		pushEOSCActions(api,
			system.NewClaimRewards(owner),
		)
	},
}

func init() {
	systemCmd.AddCommand(systemClaimRewardsCmd)
}
