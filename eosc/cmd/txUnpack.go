// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	eos "github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
)

var txUnpackCmd = &cobra.Command{
	Use:   "unpack [transaction.yaml|json]",
	Short: "Unpack a transaction produced by --output-transaction and display all its actions (for review).  This does not submit anything to the chain.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]

		cnt, err := ioutil.ReadFile(filename)
		errorCheck("reading transaction file", err)

		var tx *eos.SignedTransaction
		errorCheck("json unmarshal transaction", json.Unmarshal(cnt, &tx))

		api := getAPI()

		for _, act := range tx.ContextFreeActions {
			data, err := api.ABIBinToJSON(act.Account, eos.Name(act.Name), act.HexData)
			errorCheck("abi bin to json", err)
			act.Data = data
		}
		for _, act := range tx.Actions {
			data, err := api.ABIBinToJSON(act.Account, eos.Name(act.Name), act.HexData)
			errorCheck("abi bin to json", err)
			act.Data = data
		}

		cnt, err = json.MarshalIndent(tx, "", "  ")
		errorCheck("marshalling signed transaction", err)

		fmt.Println(string(cnt))
	},
}

func init() {
	txCmd.AddCommand(txUnpackCmd)
}
