package cmd

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getCodeCmd = &cobra.Command{
	Use:   "code [account name]",
	Short: "retrieve the code associated with an account",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		accountName := toAccount(args[0], "account name")
		codeAndABI, err := api.GetRawCodeAndABI(accountName)
		errorCheck("get code", err)

		if codeAndABI.WASMasBase64 == "" {
			errorCheck("get code", fmt.Errorf("no code has been set for account %q", accountName))
			return
		}

		normalizedWASMBase64 := codeAndABI.WASMasBase64[:len(codeAndABI.WASMasBase64)-1]
		wasm, err := base64.StdEncoding.DecodeString(normalizedWASMBase64)
		errorCheck("decode WASM base64", err)

		hash := sha256.Sum256(wasm)
		fmt.Println("Code hash: ", hex.EncodeToString(hash[:]))

		if wasmFile := viper.GetString("get-code-cmd-output-wasm"); wasmFile != "" {
			err = ioutil.WriteFile(wasmFile, wasm, 0644)
			errorCheck("writing file", err)
			fmt.Printf("Wrote WASM to %q\n", wasmFile)
		}

		if abiFile := viper.GetString("get-code-cmd-output-raw-abi"); abiFile != "" {
			if codeAndABI.ABIasBase64 != "" {
				normalizedABIBase64 := codeAndABI.ABIasBase64[:len(codeAndABI.ABIasBase64)-1]

				abi, err := base64.StdEncoding.DecodeString(normalizedABIBase64)
				errorCheck("decode ABI base64", err)
				err = ioutil.WriteFile(abiFile, abi, 0644)
				errorCheck("writing file", err)
				fmt.Printf("Wrote ABI to %q\n", abiFile)
			} else {
				errorCheck("get code", fmt.Errorf("no ABI has been set for account %q", accountName))
			}
		}

	},
}

func init() {
	getCmd.AddCommand(getCodeCmd)

	getCodeCmd.Flags().StringP("output-wasm", "", "", "Output WASM code to a file")
	getCodeCmd.Flags().StringP("output-raw-abi", "", "", "Output raw ABI to a file - If you need the JSON ABI, use `eosc get abi`")
}
