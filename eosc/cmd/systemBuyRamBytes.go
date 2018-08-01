// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"fmt"
	"strconv"

	"os"

	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var systemBuyRAMBytesCmd = &cobra.Command{
	Use:   "buyrambytes [payer] [receiver] [num bytes]",
	Short: "Buy RAM at market price, for a given number of bytes.",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		payer := toAccount(args[0], "payer")
		receiver := toAccount(args[1], "receiver")
		numBytes, err := strconv.ParseInt(args[2], 10, 64)
		errorCheck(fmt.Sprintf("invalid number of bytes %q", args[2]), err)

		if int64(uint32(numBytes)) != numBytes {
			fmt.Printf("Invalid number of bytes: capped at unsigned 32 bits.  That's probably too much RAM anyway.\n")
			os.Exit(1)
		}

		pushEOSCActions(api,
			system.NewBuyRAMBytes(payer, receiver, uint32(numBytes)),
		)
	},
}

func init() {
	systemCmd.AddCommand(systemBuyRAMBytesCmd)
}
