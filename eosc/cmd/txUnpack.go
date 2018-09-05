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

var txUnpackCmd = &cobra.Command{
	Use:   "unpack [transaction.yaml|json]",
	Short: "Unpack a transaction produced by --write-transaction and display all its actions (for review).  This does not submit anything to the chain.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]

		cnt, err := ioutil.ReadFile(filename)
		errorCheck("reading transaction file", err)

		var tx *eos.SignedTransaction
		errorCheck("json unmarshal transaction", json.Unmarshal(cnt, &tx))

		api := getAPI()

		for _, act := range tx.ContextFreeActions {
			errorCheck("context free action unpack", txUnpackAction(api, act))
		}
		for _, act := range tx.Actions {
			errorCheck("action unpack", txUnpackAction(api, act))
		}

		cnt, err = json.MarshalIndent(tx, "", "  ")
		errorCheck("marshalling signed transaction", err)

		fmt.Println(string(cnt))
	},
}

func txUnpackAction(api *eos.API, act *eos.Action) error {
	hexBytes, ok := act.Data.(string)
	if !ok {
		return fmt.Errorf("action data expected to be hex bytes as string, was %T", act.Data)
	}
	bytes, err := hex.DecodeString(hexBytes)
	if err != nil {
		return fmt.Errorf("invalid hex bytes stream: %s", err)
	}

	data, err := api.ABIBinToJSON(act.Account, eos.Name(act.Name), bytes)
	if err != nil {
		return fmt.Errorf("chain abi_bin_to_json: %s", err)
	}

	act.Data = data
	return nil
}

func init() {
	txCmd.AddCommand(txUnpackCmd)
}
