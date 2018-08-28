// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"encoding/json"

	eos "github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
)

var txCreateCmd = &cobra.Command{
	Use:   "create [contract] [action] [payload]",
	Short: "Create a transaction with a single action",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		contract := toAccount(args[0], "contract")
		action := toActionName(args[1], "action")
		payload := args[2]

		var dump map[string]interface{}
		err := json.Unmarshal([]byte(payload), &dump)
		errorCheck("[payload] is not valid JSON", err)

		api := getAPI()
		actionBinary, err := api.ABIJSONToBin(contract, eos.Name(action), dump)
		errorCheck("unable to retrieve action binary from JSON via API", err)

		pushEOSCActions(api, &eos.Action{
			Account:    contract,
			Name:       action,
			ActionData: eos.NewActionDataFromHexData([]byte(actionBinary)),
		})
	},
}

func init() {
	txCmd.AddCommand(txCreateCmd)
}
