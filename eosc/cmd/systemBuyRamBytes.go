// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"strconv"

	"os"

	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var systemBuyRamBytesCmd = &cobra.Command{
	Use:   "buyrambytes [payer] [receiver] [num bytes]",
	Short: "Buy RAM at market price, for a given number of bytes.",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		api := apiWithWallet()

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
	systemCmd.AddCommand(systemBuyRamBytesCmd)
}
