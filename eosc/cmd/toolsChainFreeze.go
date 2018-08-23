package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"encoding/hex"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/p2p"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var toolsChainFreezeCmd = &cobra.Command{
	Use:   "chain-freeze",
	Short: "Runs a p2p protocol-level proxy, and stop sync'ing the chain at the given block-num.",
	Run: func(cmd *cobra.Command, args []string) {

		chainID, err := hex.DecodeString(viper.GetString("tools-chain-freeze-cmd-chain-id"))
		errorCheck("parsing chain id", err)
		proxy := p2p.NewProxy(
			p2p.NewOutgoingPeer(viper.GetString("tools-chain-freeze-cmd-peer1-p2p-address"), chainID, "eos-proxy", false),
			p2p.NewOutgoingPeer(viper.GetString("tools-chain-freeze-cmd-peer2-p2p-address"), chainID, "eos-proxy", true),
		)

		proxy.RegisterHandler(chainFreezeHandler)
		err = proxy.Start()
		errorCheck("client start", err)
	},
}

func init() {
	toolsCmd.AddCommand(toolsChainFreezeCmd)

	toolsChainFreezeCmd.Flags().StringP("peer1-p2p-address", "", "localhost:9876", "First peer(the feed) to connect to")
	toolsChainFreezeCmd.Flags().StringP("peer2-p2p-address", "", ":19876", "Second peer(the destination) to connect to")
	toolsChainFreezeCmd.Flags().StringP("chain-id", "", "", "chain-id that will proxy")
	toolsChainFreezeCmd.Flags().IntP("on-block-modulo", "", 0, "Execute --exec-cmd each time 'block_num % module' is zero.")
	toolsChainFreezeCmd.Flags().StringP("on-actions", "", "", "Execute each time the given actions are present in a block. Format: contract1:action1,contract2:action2,...")
	toolsChainFreezeCmd.Flags().StringP("exec-cmd", "", "", "Command to execute on matching blocks")

	for _, flag := range []string{"peer1-p2p-address", "peer2-p2p-address", "exec-cmd", "on-block-modulo", "on-actions", "chain-id"} {
		if err := viper.BindPFlag("tools-chain-freeze-cmd-"+flag, toolsChainFreezeCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}

var chainFreezeHandler = p2p.HandlerFunc(func(envelope *p2p.Envelope) {
	blockModulo := viper.GetInt("tools-chain-freeze-cmd-on-block-modulo")
	actions := viper.GetString("tools-chain-freeze-cmd-on-actions")

	p2pMsg := envelope.Packet.P2PMessage
	switch m := p2pMsg.(type) {
	case *eos.SignedBlock:
		fmt.Printf("Receiving block %d sign from %s\n", m.BlockNumber(), m.Producer)

		doExec := false

		if blockModulo != 0 {
			if int(m.BlockNumber())%blockModulo == 0 {
				// run EXEC, block and continue after
				doExec = true
				goto runexec
			}
		}

		if actions != "" {
			for _, trx := range m.Transactions {
				unpacked, err := trx.Transaction.Packed.Unpack()
				if err != nil {
					fmt.Printf("Error unpacking transactions in block %d: %s\n", m.BlockNumber(), err)
					os.Exit(1)
				}

				for _, act := range unpacked.Transaction.Actions {
					actstr := fmt.Sprintf("%s:%s", act.Account, act.Name)
					if strings.Contains(actions, actstr) {
						doExec = true
						goto runexec
					}
				}
			}
		}

	runexec:
		if doExec {
			if err := runExec(); err != nil {
				fmt.Println("Error running exec:", err)
				os.Exit(1)
			}
		}
	}
})

func runExec() error {
	cmd := exec.Command(viper.GetString("tools-chain-freeze-cmd-exec-cmd"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
