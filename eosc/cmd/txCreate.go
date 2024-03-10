// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var txCreateCmd = &cobra.Command{
	Use:   "create [contract] [action] [payload]",
	Short: "Create a transaction with a single action",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		contract := toAccount(args[0], "contract")
		action := toActionName(args[1], "action")
		payload := args[2]

		forceUnique := viper.GetBool("tx-create-cmd-force-unique")

		api := getAPI()
		abi, err := api.GetABI(ctx, contract)
		errorCheck("unable to get abi", err)

		actionBinary, err := abi.ABI.EncodeAction(action, []byte(payload))
		errorCheck("unable to retrieve action binary from JSON", err)

		actions := []*eos.Action{
			&eos.Action{
				Account:    contract,
				Name:       action,
				ActionData: eos.NewActionDataFromHexData([]byte(actionBinary)),
			}}

		var contextFreeActions []*eos.Action
		if forceUnique {
			contextFreeActions = append(contextFreeActions, newNonceAction())
		}

		pushEOSCActionsAndContextFreeActions(ctx, api, contextFreeActions, actions)
	},
}

func newNonceAction() *eos.Action {
	return &eos.Action{
		Account: eos.AN("eosio.null"),
		Name:    eos.ActN("nonce"),
		ActionData: eos.NewActionData(system.Nonce{
			Value: hex.EncodeToString(generateRandomNonce()),
		}),
	}
}

func generateRandomNonce() []byte {
	// Use 48 bits of entropy to generate a valid random
	nonce := make([]byte, 6)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		errorCheck("unable to correctly generate nonce", err)
	}

	return nonce
}

func init() {
	txCmd.AddCommand(txCreateCmd)

	txCreateCmd.Flags().BoolP("force-unique", "f", false, "force the transaction to be unique. this will consume extra bandwidth and remove any protections against accidently issuing the same transaction multiple times")
}
