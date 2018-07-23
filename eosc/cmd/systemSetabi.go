package cmd

import (
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var systemSetabiCmd = &cobra.Command{
	Use:   "setabi [account name] [abi file]",
	Short: "Set ABI only on an account",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		accountName := toAccount(args[0], "account name")
		abiFile := args[1]

		action, err := system.NewSetABI(accountName, abiFile)
		errorCheck("loading abi file", err)

		pushEOSCActions(api, action)
	},
}

func init() {
	systemCmd.AddCommand(systemSetabiCmd)
}
