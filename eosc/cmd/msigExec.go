// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/eoscanada/eos-go/msig"
	"github.com/spf13/cobra"
)

// msigExecCmd represents the `eosio.msig::exec` command
var msigExecCmd = &cobra.Command{
	Use:   "exec [proposer] [proposal name] [executer]",
	Short: "Execute a transaction in the eosio.msig contract",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		proposer := toAccount(args[0], "proposer")
		proposalName := toName(args[1], "proposal name")
		executer := toAccount(args[2], "executer")

		pushEOSCActions(api,
			msig.NewExec(proposer, proposalName, executer),
		)
	},
}

func init() {
	msigCmd.AddCommand(msigExecCmd)
}
