package cmd

import (
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var systemUndelegateBWCmd = &cobra.Command{
	Use:   "undelegatebw [from] [receiver] [network bw unstake qty] [cpu bw unstake qty]",
	Short: "Undelegate some CPU and Network bandwidth.",
	Long: `Undelegate some CPU and Network bandwidth.

When undelegating bandwidth, a "refund" action will automatically be
triggered and delayed for 72 hours.  This means it takes 3 days for
you to get your EOS back and being able to transfer it. However, your
voting power is immediately altered.

See also: the "system delegatebw" command.
`,
	Args: cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		from := toAccount(args[0], "from")
		receiver := toAccount(args[1], "receiver")
		netStake, err := eos.NewEOSAssetFromString(args[2])
		errorCheck(`"network bw unstake qty" invalid`, err)
		cpuStake, err := eos.NewEOSAssetFromString(args[3])
		errorCheck(`"cpu bw unstake qty" invalid`, err)

		api := getAPI()

		pushEOSCActions(api, system.NewUndelegateBW(from, receiver, cpuStake, netStake))
	},
}

func init() {
	systemCmd.AddCommand(systemUndelegateBWCmd)
}
