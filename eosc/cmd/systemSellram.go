package cmd

import (
	"fmt"
	"strconv"

	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var systemSellRAMCmd = &cobra.Command{
	Use:   "sellram [account_name] [num bytes]",
	Short: "Sell the [num bytes] amount of bytes of RAM on the RAM market.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		accountName := toAccount(args[0], "account name")
		numBytes, err := strconv.ParseInt(args[1], 10, 64)
		errorCheck(fmt.Sprintf("invalid number of bytes %q", args[1]), err)

		pushEOSCActions(api,
			system.NewSellRAM(accountName, uint64(numBytes)),
		)
	},
}

func init() {
	systemCmd.AddCommand(systemSellRAMCmd)
}
