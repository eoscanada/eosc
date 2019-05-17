// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	eos "github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

var txPushCmd = &cobra.Command{
	Use:   "push [transaction.json | string]",
	Short: "Push a signed transaction to the chain. Must be done online.",
	Long: `
The first argument interpretation is heuristics based, and depends on its
content.

If it's an existing file on disk, the first argument is assumed to be a
filename containing a transaction in JSON format and is used right away.

Otherwise, if the first argument looks like a JSON string (basic checks
just enough to avoid ambiguity with a plain filename), then the input is
considered a JSON transaction.

Otherwise, if the argument is not deem being a valid JSON input, it's
assumed to be a filename that is not present and a proper error message
is displayed.
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cnt, err := readTransaction(args[0])
		errorCheck("reading transaction", err)

		chainID := gjson.GetBytes(cnt, "chain_id").String()
		hexChainID, _ := hex.DecodeString(chainID)

		var signedTx *eos.SignedTransaction
		errorCheck("json unmarshal transaction", json.Unmarshal(cnt, &signedTx))

		api := getAPI()

		packedTx, err := signedTx.Pack(eos.CompressionNone)
		errorCheck("packing transaction", err)

		pushTransaction(api, packedTx, eos.SHA256Bytes(hexChainID))
	},
}

func readTransaction(argument string) ([]byte, error) {
	fileInfo, err := os.Stat(argument)
	if err == nil {
		if fileInfo.IsDir() {
			return nil, fmt.Errorf("argument %q is a directory", argument)
		}

		return ioutil.ReadFile(argument)
	}

	if looksLikeTransactionJSON(argument) {
		return []byte(argument), nil
	}

	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file %q does not exist", argument)
	}

	return nil, fmt.Errorf("unable to check file %q: %s", argument, err)
}

var jsonFieldRegexp = regexp.MustCompile(`".+"\s*:\s*".*"`)

func looksLikeTransactionJSON(input string) bool {
	return strings.Contains(input, "{") && jsonFieldRegexp.MatchString(input)
}

func init() {
	txCmd.AddCommand(txPushCmd)
}
