// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getTableCmd = &cobra.Command{
	Use:   "table [contract] [scope] [table]",
	Short: "Fetch data from a table in a contract on chain",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		response, err := api.GetTableRows(
			eos.GetTableRowsRequest{
				Code:       args[0],
				Scope:      args[1],
				Table:      args[2],
				LowerBound: viper.GetString("get-table-cmd-lower-bound"),
				UpperBound: viper.GetString("get-table-cmd-upper-bound"),
				Limit:      uint32(viper.GetInt("get-table-cmd-limit")),
				KeyType:    viper.GetString("get-table-cmd-key-type"),
				Index:      viper.GetString("get-table-cmd-index"),
				EncodeType: viper.GetString("get-table-cmd-encode-type"),
				JSON:       !(viper.GetBool("get-table-cmd-output-binary")),
			},
		)
		errorCheck("get table rows", err)

		data, err := json.MarshalIndent(response, "", "  ")
		errorCheck("json marshal", err)

		fmt.Println(string(data))
	},
}

func init() {
	getCmd.AddCommand(getTableCmd)

	getTableCmd.Flags().String("lower-bound", "", "Lower bound (incluse) value of key, defaults to first.")
	getTableCmd.Flags().String("upper-bound", "", "Upper bound (exclusive) value of key, defaults to first.")
	getTableCmd.Flags().String("key-type", "", "The key type of --index, primary only supports (i64), all others support (i64, i128, i256, float64, float128, ripemd160, sha256). Special type 'name' indicates an account name.")
	getTableCmd.Flags().String("index", "", "Index number, 1 - primary (first), 2 - secondary index (in order defined by multi_index), 3 - third index, etc. Number or name of index can be specified, e.g. 'secondary' or '2'.")
	getTableCmd.Flags().String("encode-type", "", "The encoding type of key-type (i64 , i128 , float64, float128) only support decimal encoding e.g. 'dec.  i256 - supports both 'dec' and 'hex', ripemd160 and sha256 is 'hex' only.")
	getTableCmd.Flags().Int("limit", 100, "Maximum number of rows to return.")
	getTableCmd.Flags().Bool("output-binary", false, "Outputs the row-level data as hex-encoded binary instead of deserializing using the ABI")
}
