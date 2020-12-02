// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"context"
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
		ctx := context.Background()

		var tx *eos.SignedTransaction
		rawTx, err := transactionFromFileOrCLI(args[0], &tx)
		errorCheck("reading transaction", err)

		api := getAPI()

		for _, act := range tx.ContextFreeActions {
			errorCheck("context free action unpack", txUnpackAction(ctx, api, act))
		}
		for _, act := range tx.Actions {
			errorCheck("action unpack", txUnpackAction(ctx, api, act))
		}

		rawTx, err = json.MarshalIndent(tx, "", "  ")
		errorCheck("marshalling signed transaction", err)

		fmt.Println(string(rawTx))
	},
}

func txUnpackAction(ctx context.Context, api *eos.API, act *eos.Action) error {
	hexBytes, ok := act.Data.(string)
	if !ok {
		return fmt.Errorf("action data expected to be hex bytes as string, was %T", act.Data)
	}
	bytes, err := hex.DecodeString(hexBytes)
	if err != nil {
		return fmt.Errorf("invalid hex bytes stream: %s", err)
	}

	data, err := api.ABIBinToJSON(ctx, act.Account, eos.Name(act.Name), bytes)
	if err != nil {
		return fmt.Errorf("chain abi_bin_to_json: %s", err)
	}

	act.Data = data
	return nil
}

func init() {
	txCmd.AddCommand(txUnpackCmd)
}

func transactionFromFileOrCLI(input string, into interface{}) (cnt []byte, err error) {
	var filename string
	if input[0] == '{' {
		cnt = []byte(input)
		filename = "stdin"
	} else {
		cnt, err = ioutil.ReadFile(input)
		if err != nil {
			return nil, fmt.Errorf("reading transaction file %q: %w", input, err)
		}
		filename = input
	}

	err = json.Unmarshal(cnt, into)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal transaction from %s: %w", filename, err)
	}

	return
}
