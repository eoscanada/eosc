// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	eos "github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var txSignCmd = &cobra.Command{
	Use:   "sign [transaction.yaml|json]",
	Short: "Sign a transaction produced by --write-transaction and submit it to the chain (unless --write-transaction is passed again).",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]

		cnt, err := ioutil.ReadFile(filename)
		errorCheck("reading transaction file", err)

		var tx *eos.Transaction
		errorCheck("json unmarshal transaction", json.Unmarshal(cnt, &tx))

		api := getAPI()

		var chainID eos.SHA256Bytes
		if infileChainID := gjson.Get(string(cnt), "chain_id").String(); infileChainID != "" {
			chainID = toSHA256Bytes(infileChainID, fmt.Sprintf("chain_id field in %q", filename))
		} else if cliChainID := viper.GetString("global-offline-chain-id"); cliChainID != "" {
			chainID = toSHA256Bytes(cliChainID, "--offline-chain-id")
		} else {
			// getInfo
			resp, err := api.GetInfo()
			errorCheck("get info", err)
			chainID = resp.ChainID
		}

		signedTx, packedTx := optionallySignTransaction(tx, chainID, api)

		optionallyPushTransaction(signedTx, packedTx, chainID, api)
	},
}

func init() {
	txCmd.AddCommand(txSignCmd)

	// txSignCmd.Flags().StringP("hash", "", "", "Hash of the proposition, as defined by the proposition itself")

	// for _, flag := range []string{"hash"} {
	// 	if err := viper.BindPFlag("tx-sign-cmd-"+flag, txSignCmd.Flags().Lookup(flag)); err != nil {
	// 		panic(err)
	// 	}
	// }
}
