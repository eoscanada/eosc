package cmd

import (
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var systemUnregisterProducerCmd = &cobra.Command{
	Use:   "unregprod [account_name]",
	Short: "Unregister producer account temporarily.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		accountName := toAccount(args[0], "account name")

		pushEOSCActions(api,
			system.NewUnregProducer(accountName),
		)
	},
}

func init() {
	systemCmd.AddCommand(systemUnregisterProducerCmd)
}
