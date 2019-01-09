// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"

	eos "github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
)

var txIDCmd = &cobra.Command{
	Use:   "id [transaction.json]",
	Short: "Print the transaction ID for a given transaction file.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]

		cnt, err := ioutil.ReadFile(filename)
		errorCheck("reading file", err)

		var stx *eos.SignedTransaction
		errorCheck("parsing JSON content", json.Unmarshal(cnt, &stx))

		ptx, err := stx.Pack(eos.CompressionNone)
		errorCheck("packing transaction", err)

		ptxID, err := ptx.ID()
		errorCheck("transaction ID", err)

		fmt.Println(hex.EncodeToString(ptxID))
	},
}

func init() {
	txCmd.AddCommand(txIDCmd)
}
