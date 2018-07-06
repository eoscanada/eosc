package cmd

import (
	"fmt"

	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var voteProxyCmd = &cobra.Command{
	Use:   "proxy [voter name] [proxy name]",
	Short: "Cast your vote for a proxy voter",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		voterName := toAccount(args[0], "voter name")
		proxyName := toAccount(args[1], "proxy name")

		fmt.Printf("Voter [%s] voting for proxy: %s\n", voterName, proxyName)

		pushEOSCActions(api,
			system.NewVoteProducer(
				voterName,
				proxyName,
			),
		)
	},
}

func init() {
	voteCmd.AddCommand(voteProxyCmd)
}
