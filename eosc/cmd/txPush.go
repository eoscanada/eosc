// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"encoding/json"
	"io/ioutil"

	eos "github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
)

var txPushCmd = &cobra.Command{
	Use:   "push [transaction.yaml|json]",
	Short: "Push a signed transaction to the chain.  Must be done online.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]

		cnt, err := ioutil.ReadFile(filename)
		errorCheck("reading transaction file", err)

		var signedTx *eos.SignedTransaction
		errorCheck("json unmarshal transaction", json.Unmarshal(cnt, &signedTx))

		api := getAPI()

		packedTx, err := signedTx.Pack(eos.CompressionNone)
		errorCheck("packing transaction", err)

		pushTransaction(api, packedTx)
	},
}

func init() {
	txCmd.AddCommand(txPushCmd)

	// txPushCmd.Flags().StringP("hash", "", "", "Hash of the proposition, as defined by the proposition itself")

	// for _, flag := range []string{"hash"} {
	// 	if err := viper.BindPFlag("tx-push-cmd-"+flag, txPushCmd.Flags().Lookup(flag)); err != nil {
	// 		panic(err)
	// 	}
	// }
}
