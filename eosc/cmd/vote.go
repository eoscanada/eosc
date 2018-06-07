package cmd

import (
	"fmt"
	"os"

	"sort"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var voteCmd = &cobra.Command{
	Use:   "vote",
	Short: "Command to vote for block producers or proxy",
}
var voteProducersCmd = &cobra.Command{
	Use:   "producers [voter name] [producer list]",
	Short: "Cast your vote for 1 to 30 producers. View them with 'list'",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var producerNames = make([]eos.AccountName, 0, 0)

		producerStringNames := args[1:]
		sort.Strings(producerStringNames)
		for _, producerString := range producerStringNames {
			producerNames = append(producerNames, eos.AccountName(producerString))
		}

		if len(producerNames) == 0 {
			fmt.Printf("No producer provided")
			os.Exit(1)
		}

		api, err := apiWithWallet()
		if err != nil {
			fmt.Printf("Error initiating api, %s\n", err.Error())
			os.Exit(1)
		}

		voterName := args[0]

		fmt.Printf("Voter [%s] voting for: %s\n", voterName, producerNames)
		_, err = api.SignPushActions(
			system.NewVoteProducer(
				eos.AccountName(voterName),
				"",
				producerNames...,
			),
		)

		if err != nil {
			fmt.Printf("Error during vote, %s\n", err.Error())
			os.Exit(1)
		}

		fmt.Println("Vote sent to chain.")
	},
}

var voteProxyCmd = &cobra.Command{
	Use:   "proxy [voter name] [proxy name]",
	Short: "Cast your vote for a proxy voter",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		proxyName := eos.AccountName(args[1])

		api, err := apiWithWallet()
		if err != nil {
			fmt.Printf("Error initiating api, %s\n", err.Error())
			os.Exit(1)
		}

		voterName := args[0]

		fmt.Printf("Voter [%s] voting for proxy: %s\n", voterName, proxyName)
		_, err = api.SignPushActions(
			system.NewVoteProducer(
				eos.AccountName(voterName),
				proxyName,
				[]eos.AccountName{}...,
			),
		)

		if err != nil {
			fmt.Printf("Error during vote, %s\n", err.Error())
			os.Exit(1)
		}

		fmt.Println("Vote sent to chain.")
	},
}

func init() {
	RootCmd.AddCommand(voteCmd)
	voteCmd.AddCommand(voteProducersCmd)
	voteCmd.AddCommand(voteProxyCmd)
}
