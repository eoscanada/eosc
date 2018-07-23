package cmd

import (
	"fmt"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var systemDelegateBWCmd = &cobra.Command{
	Use:   "delegatebw [from_account_name] [receiver_account_name]",
	Short: "Delegate cpu and network bandwidth to receiver",
	Long:  "Delegate cpu and network bandwidth to receiver",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		from := eos.AccountName(args[0])
		receiver := eos.AccountName(args[1])

		cpuStakeStr := viper.GetString("system-delegatebw-cmd-stake-cpu")
		netStakeStr := viper.GetString("system-delegatebw-cmd-stake-net")

		if cpuStakeStr == "" {
			errorCheck("missing argument", fmt.Errorf("--stake-cpu missing"))
		}
		if netStakeStr == "" {
			errorCheck("missing argument", fmt.Errorf("--stake-net missing"))
		}

		cpuStake, err := eos.NewEOSAssetFromString(cpuStakeStr)
		errorCheck("--stake-cpu invalid", err)
		netStake, err := eos.NewEOSAssetFromString(netStakeStr)
		errorCheck("--stake-net invalid", err)

		doTransfer := viper.GetBool("system-delegatebw-cmd-transfer")

		api := apiWithWallet()

		pushEOSCActions(api,
			system.NewDelegateBW(from, receiver, cpuStake, netStake, doTransfer),
		)
	},
}

func init() {
	systemCmd.AddCommand(systemDelegateBWCmd)
	systemDelegateBWCmd.Flags().StringP("stake-cpu", "", "", "Amount of EOS to stake for CPU bandwidth (required)")
	systemDelegateBWCmd.Flags().StringP("stake-net", "", "", "Amount of EOS to stake for Network bandwidth (required)")
	systemDelegateBWCmd.Flags().BoolP("transfer", "", false, "Transfer voting power and right to unstake EOS to receiver")

	for _, flag := range []string{"stake-cpu", "stake-net", "transfer"} {
		if err := viper.BindPFlag("system-delegatebw-cmd-"+flag, systemDelegateBWCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}
