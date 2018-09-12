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

		if codeAndABI.AccountName == accountName {
			//TODO: Ugly hack !?! -- does not work when string len 0
			fixedWASMBase64 := codeAndABI.WASMasBase64[:len(codeAndABI.WASMasBase64)-1]
			fixedABIBase64 := codeAndABI.ABIasBase64[:len(codeAndABI.ABIasBase64)-1]

			bytecode, err := base64.StdEncoding.DecodeString(fixedWASMBase64)
			errorCheck("decode WASM base64", err)

			hash := sha256.Sum256(bytecode)
			fmt.Println("Code hash:", hex.EncodeToString(hash[:]))

			if wasmFile := viper.GetString("get-code-cmd-output-wasm"); wasmFile != "" {
				err = ioutil.WriteFile(wasmFile, bytecode, 0644)
				errorCheck("writing file", err)

				fmt.Printf("Wrote WASM to %q\n", wasmFile)
			}

			if abiFile := viper.GetString("get-code-cmd-output-abi"); abiFile != "" {
				abi, err := base64.StdEncoding.DecodeString(fixedABIBase64)
				errorCheck("decode ABI base64", err)
				err = ioutil.WriteFile(abiFile, abi, 0644)
				errorCheck("writing file", err)

				fmt.Printf("Wrote ABI to %q\n", abiFile)

				// data, err := json.MarshalIndent(codeAndABI.ABI, "", "  ")
				// errorCheck("json marshal", err)
				// fmt.Println(string(data))

				// fmt.Printf("Wrote ABI to %q\n", abiFile)
			}
		} else {
			errorCheck("get code", fmt.Errorf("unable to find account %s", accountName))
		}
	},
}

func init() {
	getCmd.AddCommand(getCodeCmd)

	getCodeCmd.Flags().StringP("output-wasm", "", "", "Output WASM code to a file")
	getCodeCmd.Flags().StringP("output-abi", "", "", "Output ABI to a file")

	for _, flag := range []string{"output-wasm", "output-abi"} {
		if err := viper.BindPFlag("get-code-cmd-"+flag, getCodeCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}

}
