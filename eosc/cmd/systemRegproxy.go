package cmd

import (
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var systemRegisterProxyCmd = &cobra.Command{
	Use:   "regproxy [account_name]",
	Short: "Register an account as a voting proxy",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()
		accountName := toAccount(args[0], "account name")
		pushEOSCActions(api,
			system.NewRegProxy(accountName, true),
		)
	},
}

func init() {
	systemCmd.AddCommand(systemRegisterProxyCmd)
}
