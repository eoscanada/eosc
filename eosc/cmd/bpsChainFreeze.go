package cmd

import (
	"fmt"
	"os"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/p2p"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var bpsChainFreezeCmd = &cobra.Command{
	Use:   "chain-freeze",
	Short: "Freeze the chain by proxying p2p blocks until a block including updateauth actions is passed through, then block/shutdown.",
	Run: func(cmd *cobra.Command, args []string) {

		proxy := p2p.Proxy{
			Routes: []*p2p.Route{
				{From: viper.GetString("listening-address"), To: viper.GetString("target-p2p-address")},
			},
			Handlers: []p2p.Handler{chainFreezeHandler},
		}

		proxy.Start()

	},
}

func init() {
	bpsCmd.AddCommand(bpsChainFreezeCmd)

	bpsChainFreezeCmd.Flags().StringP("target-p2p-address", "t", "localhost:9876", "return producers info in json")
	bpsChainFreezeCmd.Flags().StringP("listening-address", "", ":19876", "return producers info in json")

	for _, flag := range []string{"target-p2p-address", "listening-address"} {
		if err := viper.BindPFlag(flag, bpsChainFreezeCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}

var chainFreezeHandler = p2p.HandlerFunc(func(msg p2p.Message) {
	p2pMsg := msg.Envelope.P2PMessage
	switch m := p2pMsg.(type) {
	case *eos.SignedBlock:
		fmt.Printf("Receiving block %d sign from %s\n", m.BlockNumber(), m.Producer)
		for _, tx := range m.Transactions {
			signTransaction, err := tx.Transaction.Packed.Unpack()
			if err != nil {
				fmt.Println("Error: unpack, ", err.Error())
			}
			for _, action := range signTransaction.Actions {
				fmt.Printf("\tReceived action %s::%s\n", action.Account, action.Name)

				// ACTION on which to close connection.
				if action.Name == "updateauth" {
					fmt.Println("Closing connection, enjoy your frozen chain.")
					os.Exit(0)
				}
			}
		}
	default:
		//fmt.Println("found type: Default")
	}
})
