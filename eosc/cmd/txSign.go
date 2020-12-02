// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"context"
	"fmt"

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
		ctx := context.Background()

		var tx *eos.Transaction
		rawTx, err := transactionFromFileOrCLI(args[0], &tx)
		errorCheck("reading transaction", err)

		api := getAPI()

		var chainID eos.SHA256Bytes
		if infileChainID := gjson.Get(string(rawTx), "chain_id").String(); infileChainID != "" {
			chainID = toSHA256Bytes(infileChainID, fmt.Sprintf("chain_id field in source transaction"))
		} else if cliChainID := viper.GetString("global-offline-chain-id"); cliChainID != "" {
			chainID = toSHA256Bytes(cliChainID, "--offline-chain-id")
		} else {
			// getInfo
			resp, err := api.GetInfo(ctx)
			errorCheck("get info", err)
			chainID = resp.ChainID
		}

		signedTx, packedTx := optionallySignTransaction(ctx, tx, chainID, api, false)

		optionallyPushTransaction(ctx, signedTx, packedTx, chainID, api)
	},
}

func init() {
	txCmd.AddCommand(txSignCmd)

	// txSignCmd.Flags().StringP("hash", "", "", "Hash of the proposition, as defined by the proposition itself")
}
