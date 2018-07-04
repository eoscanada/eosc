package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getCodeCmd = &cobra.Command{
	Use:   "code [account name]",
	Short: "retrieve account information for a given name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		accountName := toAccount(args[0], "account name")
		code, err := api.GetCode(accountName)
		errorCheck("get code", err)

		fmt.Println("Code hash:", code.CodeHash)

		outFile := viper.GetString("get-code-cmd-output")
		if outFile != "" {
			// fmt.Println("MAMA", code.WASM)
			// bytecode, err := hex.DecodeString(code.WASM)
			// errorCheck("decode wasm hex", err)
			// err = ioutil.WriteFile(outFile, code.WASM, 0644)
			// errorCheck("writing file", err)
			// fmt.Printf("Wrote wasm to %q\n", outFile)
		}
	},
}

func init() {
	getCmd.AddCommand(getCodeCmd)

	getCodeCmd.Flags().StringP("output", "", "", "Output .wasm code to filename")

	for _, flag := range []string{"output"} {
		if err := viper.BindPFlag("get-code-cmd-"+flag, getCodeCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}

}
