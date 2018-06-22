package cmd

import (
	"fmt"
	"os"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/p2p"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var toolsChainFreezeCmd = &cobra.Command{
	Use:   "chain-freeze",
	Short: "Freeze the chain by proxying p2p blocks until a block including updateauth actions is passed through, then block/shutdown.",
	Run: func(cmd *cobra.Command, args []string) {

		proxy := p2p.Proxy{
			Routes: []*p2p.Route{
				{From: viper.GetString("tools-chain-freeze-cmd-listening-address"), To: viper.GetString("tools-chains-freeze-cmd-target-p2p-address")},
			},
			Handlers: []p2p.Handler{chainFreezeHandler},
		}

		proxy.Start()

	},
}

func init() {
	toolsCmd.AddCommand(toolsChainFreezeCmd)

	toolsChainFreezeCmd.Flags().StringP("target-p2p-address", "t", "localhost:9876", "return producers info in json")
	toolsChainFreezeCmd.Flags().StringP("listening-address", "", ":19876", "return producers info in json")
	toolsChainFreezeCmd.Flags().IntP("block-num", "n", 0, "Last block to let through before exiting.")

	for _, flag := range []string{"target-p2p-address", "listening-address", "block-num"} {
		if err := viper.BindPFlag("tools-chain-freeze-cmd-"+flag, toolsChainFreezeCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}

var chainFreezeHandler = p2p.HandlerFunc(func(msg p2p.Message) {
	maxBlock := viper.GetInt("tools-chain-freeze-cmd-block-num")

	p2pMsg := msg.Envelope.P2PMessage
	switch m := p2pMsg.(type) {
	case *eos.SignedBlock:
		fmt.Printf("Receiving block %d sign from %s\n", m.BlockNumber(), m.Producer)
		if m.BlockNumber() >= uint32(maxBlock) {
			fmt.Println("Closing connection, enjoy your frozen chain.")
			os.Exit(0)
		}
	}
})
