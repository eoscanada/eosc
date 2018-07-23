package cmd

import (
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var systemSetcontractCmd = &cobra.Command{
	Use:   "setcontract [account name] [wasm file] [abi file]",
	Short: "Set both code and ABI on an account",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		accountName := toAccount(args[0], "account name")
		wasmFile := args[1]
		abiFile := args[2]

		actions, err := system.NewSetContract(accountName, wasmFile, abiFile)
		errorCheck("loading files", err)

		pushEOSCActions(api,
			actions...,
		)
	},
}

func init() {
	systemCmd.AddCommand(systemSetcontractCmd)
}
