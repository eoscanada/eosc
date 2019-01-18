package cmd

import (
	"encoding/json"
	"io/ioutil"

	eos "github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var toolsProducerJSONCmd = &cobra.Command{
	Use:   "producerjson [account] [file.json]",
	Short: "Publish a producer json file to a producerjson-compatible contract.",
	Long: `Publish a producer json file to a producerjson-compatible contract.

Reference: https://github.com/greymass/producerjson
`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		producerAccount := toAccount(args[0], "sold account")
		fileName := args[1]

		cnt, err := ioutil.ReadFile(fileName)
		errorCheck("reading json file", err)

		// TODO: eventually do some validation on the producerjson content
		var packme map[string]interface{}
		errorCheck("file contains invalid json", json.Unmarshal(cnt, &packme))
		packedCnt, err := json.Marshal(packme)
		errorCheck("packing json more tightly", err)

		type producerJSONSet struct {
			Owner eos.AccountName `json:"owner"`
			JSON  string          `json:"json"`
		}
		pushEOSCActions(api, &eos.Action{
			Account: eos.AccountName(viper.GetString("tools-producerjson-cmd-target-contract")),
			Name:    eos.ActionName("set"),
			Authorization: []eos.PermissionLevel{
				{Actor: producerAccount, Permission: eos.PermissionName("active")},
			},
			ActionData: eos.NewActionData(producerJSONSet{
				Owner: producerAccount,
				JSON:  string(packedCnt),
			}),
		})

	},
}

func init() {
	toolsCmd.AddCommand(toolsProducerJSONCmd)

	toolsProducerJSONCmd.Flags().StringP("target-contract", "", "producerjson", "Target producerjson contract")
}
