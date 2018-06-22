package cmd

import (
	"fmt"
	"os"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var bidNewNameCmd = &cobra.Command{
	Use:   "bidname [bidder] [newname] [bid]",
	Short: "Command to name bidding",
	Long: `Command to name bidding.

bidder TEXT   The bidding account (required)
newname TEXT  The bidding name (required)
bid TEXT   The amount of EOS to bid (required)

You can bidding a new name with:

    eosc bidname myname eos "10.0000 EOS"
`,
	Args: cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		bidder := eos.AccountName(args[0])
		newname := eos.AccountName(args[1])
		bid, err := eos.NewAsset(args[2])
		if err != nil {
			fmt.Printf("Error get bid, %s\n", err.Error())
			os.Exit(1)
		}

		api, err := apiWithWallet()
		if err != nil {
			fmt.Printf("Error initiating api, %s\n", err.Error())
			os.Exit(1)
		}

		fmt.Printf("[%s] bidding for: %s , amount=%d precision=%d symbol=%s\n", bidder, newname, bid.Amount, bid.Symbol.Precision, bid.Symbol.Symbol)
		_, err = api.SignPushActions(
			system.NewBidname(
				bidder,
				newname,
				bid,
			),
		)

		if err != nil {
			fmt.Printf("Error during bidding, %s\n", err.Error())
			os.Exit(1)
		}

		fmt.Println("Bidname sent to chain.")
	},
}

func init() {
	RootCmd.AddCommand(bidNewNameCmd)
}
