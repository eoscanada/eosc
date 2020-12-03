// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"context"
	"fmt"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
	"github.com/eoscanada/eosc/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	systemCmd.AddCommand(systemNewAccountCmd)

	systemNewAccountCmd.Flags().String("stake-cpu", "", "Amount of EOS to stake for CPU bandwidth (required)")
	systemNewAccountCmd.Flags().String("stake-net", "", "Amount of EOS to stake for Network bandwidth (required)")
	systemNewAccountCmd.Flags().Bool("transfer", false, "Transfer voting power and right to unstake EOS to receiver")

	systemNewAccountCmd.Flags().Int("buy-ram-kbytes", 8, "The amount of RAM kibibytes (KiB) to purchase for the new account.  Defaults to 8 KiB.")
	systemNewAccountCmd.Flags().String("buy-ram", "", "The amount of EOS to spend to buy RAM for the new account (at current EOS/RAM market price)")

	systemNewAccountCmd.Flags().StringSlice("additional-actions", []string{"delegatebw", "buyram"}, "Action types to include in the generated transactions. Defaults to EOS Mainnet behavior. Options: 'delegatebw' (uses --stake-cpu, --stake-net and --transfer), 'buyram' (uses --buy-ram-kbytes or --buy-ram), 'setpriv' (makes account privileged, only possible if creator is 'eosio')")
}

var systemNewAccountCmd = &cobra.Command{
	Use:   "newaccount [creator] [new_account_name] [owner authority] [active authority]",
	Short: "Create a new account.",
	Long: `Create a new account

Both [owner authority] and [active authority] are expressed using this short-form syntax:
` + shortFormAuthHelp,
	Args: cobra.ExactArgs(4),
	Run:  systemNewaccountRun,
}

func systemNewaccountRun(cmd *cobra.Command, args []string) {
	creator := toAccount(args[0], "creator")
	newAccount := toAccount(args[1], "new account name")

	var actions []*eos.Action
	var ownerAuth *eos.Authority
	var activeAuth *eos.Authority
	var err error

	ownerAuth, err = cli.ParseShortFormAuth(args[2])
	errorCheck("parsing owner auth", err)
	errorCheck("invalid owner auth", ValidateAuth(ownerAuth))

	activeAuth, err = cli.ParseShortFormAuth(args[3])
	errorCheck("parsing active auth", err)
	errorCheck("invalid active auth", ValidateAuth(activeAuth))

	actions = append(actions, system.NewCustomNewAccount(creator, newAccount, *ownerAuth, *activeAuth))

	addActions := map[string]bool{}
	for _, act := range viper.GetStringSlice("system-newaccount-cmd-additional-actions") {
		if !(map[string]bool{"delegatebw": true, "buyram": true, "setpriv": true}[act]) {
			errorCheck("invalid additional-actions", fmt.Errorf("%q is not a valid action type", act))
		}
		addActions[act] = true
	}

	if addActions["delegatebw"] {
		cpuStakeStr := viper.GetString("system-newaccount-cmd-stake-cpu")
		netStakeStr := viper.GetString("system-newaccount-cmd-stake-net")

		if cpuStakeStr == "" {
			errorCheck("missing argument", fmt.Errorf("--stake-cpu missing"))
		}
		if netStakeStr == "" {
			errorCheck("missing argument", fmt.Errorf("--stake-net missing"))
		}

		cpuStake := toCoreAsset(cpuStakeStr, "--stake-cpu")
		netStake := toCoreAsset(netStakeStr, "--stake-net")

		doTransfer := viper.GetBool("system-newaccount-cmd-transfer")
		if cpuStake.Amount != 0 || netStake.Amount != 0 {
			actions = append(actions, system.NewDelegateBW(creator, newAccount, cpuStake, netStake, doTransfer))
		} else if doTransfer {
			errorCheck("--transfer invalid", fmt.Errorf("nothing was staked, so nothing to transfer"))
		}
	}

	if addActions["buyram"] {
		buyRAM := viper.GetString("system-newaccount-cmd-buy-ram")
		if buyRAM != "" {
			buyRAMAmount := toCoreAsset(buyRAM, "--buy-ram")
			actions = append(actions, system.NewBuyRAM(creator, newAccount, uint64(buyRAMAmount.Amount)))
		} else {
			buyRAMBytes := viper.GetInt("system-newaccount-cmd-buy-ram-kbytes")
			actions = append(actions, system.NewBuyRAMBytes(creator, newAccount, uint32(buyRAMBytes*1024)))
		}
	}

	if addActions["setpriv"] {
		actions = append(actions, system.NewSetPriv(newAccount))
	}

	api := getAPI()

	pushEOSCActions(context.Background(), api, actions...)
}
