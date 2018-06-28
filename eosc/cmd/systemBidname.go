package cmd

import (
	"fmt"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var bidNewNameCmd = &cobra.Command{
	Use:   "bidname [bidder_account_name] [newname] [bid_asset]",
	Short: "Bidding a name with asset.",
	Long: `Command to name bidding.

bidder_account_name TEXT  The bidding account name(required)
newname TEXT  The bidding name (required)
bid_asset TEXT  The amount of EOS to bid (required)

You can bidding a new name with:

    eosc system bidname your_account_name eos "10.0000 EOS"
`,
	Args: cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		api := apiWithWallet()

		bidder := eos.AccountName(args[0])
		newname := eos.AccountName(args[1])
		bidAsset, err := eos.NewAsset(args[2])
		errorCheck("Error get bid", err)

		fmt.Printf("[%s] bidding for: %s , amount=%d precision=%d symbol=%s\n", bidder, newname, bidAsset.Amount, bidAsset.Symbol.Precision, bidAsset.Symbol.Symbol)

		pushEOSCActions(api,
			system.NewBidname(bidder, newname, bidAsset),
		)
	},
}

func init() {
	systemCmd.AddCommand(bidNewNameCmd)
}
