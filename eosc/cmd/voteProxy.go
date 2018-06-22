package cmd

import (
	"fmt"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var voteProxyCmd = &cobra.Command{
	Use:   "proxy [voter name] [proxy name]",
	Short: "Cast your vote for a proxy voter",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		api := apiWithWallet()

		voterName := args[0]
		proxyName := eos.AccountName(args[1])

		fmt.Printf("Voter [%s] voting for proxy: %s\n", voterName, proxyName)

		pushEOSCActions(api,
			system.NewVoteProducer(
				eos.AccountName(voterName),
				proxyName,
			),
		)
	},
}

func init() {
	voteCmd.AddCommand(voteProxyCmd)
}
