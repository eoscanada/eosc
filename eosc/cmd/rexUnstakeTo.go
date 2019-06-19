// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/rex"
	"github.com/spf13/cobra"
)

var rexUnstakeTo = &cobra.Command{
	Use:   "unstake-to [staker] [staked to] [net] [cpu]",
	Short: "Utilize staked tokens to purchase REX tokens.",
	Long:  "Use this action to utilize staked tokens to purchase REX tokens. [staker] is the owner of the tokens, [staked to] is the account who has access to those resources.",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		staker := toAccount(args[0], "staker")
		stakedTo := toAccount(args[1], "staked to")
		net := toCoreAsset(args[2], "net")
		cpu := toCoreAsset(args[3], "cpu")

		pushEOSCActions(getAPI(), rex.NewUnstakeToREX(
			staker,
			stakedTo,
			net,
			cpu,
		))
	},
}

func init() {
	rexCmd.AddCommand(rexUnstakeTo)
}
