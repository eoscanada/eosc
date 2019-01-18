package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ryanuber/columnize"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eosc/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getAccountCmd = &cobra.Command{
	Use:   "account [account name]",
	Short: "retrieve account information for a given name",
	Long:  "retrieve account information for a given name.  For a json dump, append the argument --json.",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		accountName := toAccount(args[0], "account name")
		account, err := api.GetAccount(accountName)
		errorCheck("get account", err)

		if viper.GetBool("get-account-cmd-json") == true {
			data, err := json.MarshalIndent(account, "", "  ")
			errorCheck("json marshal", err)
			fmt.Println(string(data))
			return
		}
		printAccount(account)
	},
}

func printAccount(account *eos.AccountResp) {
	if account != nil {
		// dereference this so we can safely mutate it to accomodate uninitialized symbols
		act := *account
		if act.SelfDelegatedBandwidth.CPUWeight.Symbol.Symbol == "" {
			act.SelfDelegatedBandwidth.CPUWeight.Symbol = act.TotalResources.CPUWeight.Symbol
		}
		if act.SelfDelegatedBandwidth.NetWeight.Symbol.Symbol == "" {
			act.SelfDelegatedBandwidth.NetWeight.Symbol = act.TotalResources.CPUWeight.Symbol
		}
		cfg := &columnize.Config{
			NoTrim: true,
		}

		for _, s := range []string{
			cli.FormatBasicAccountInfo(&act, cfg),
			cli.FormatPermissions(&act, cfg),
			cli.FormatMemory(&act, cfg),
			cli.FormatNetworkBandwidth(&act, cfg),
			cli.FormatCPUBandwidth(&act, cfg),
			cli.FormatBalances(&act, cfg),
			cli.FormatProducers(&act, cfg),
			cli.FormatVoterInfo(&act, cfg),
		} {
			fmt.Println(s)
			fmt.Println("")
		}
	}
}

func init() {
	getCmd.AddCommand(getAccountCmd)
	getAccountCmd.Flags().BoolP("json", "", false, "pass if you wish to see account printed as json")
}
