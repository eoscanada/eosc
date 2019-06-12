package cmd

import (
	"path/filepath"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eosc/bios"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var bootCmd = &cobra.Command{
	Use:   "boot [boot_sequence.yaml]",
	Short: "Boot a fresh network, using the now famous eos-bios.",
	Long: `Boot a fresh network, using the now famous eos-bios.

Use one of the boot sequences in https://github.com/eoscanada/eosc/tree/master/bootseqs
to setup a clean EOSIO blockchain, with the features you like.

Use a base config over there, run your node, create a new Vault and use it
to bootstrap your chain by running 'eosc boot'.
`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()
		attachWallet(api)

		logger := bios.NewLogger()
		logger.Debug = viper.GetBool("global-debug")

		b := bios.NewBIOS(logger, viper.GetString("boot-cmd-cache-path"), api)
		b.WriteActions = viper.GetBool("boot-cmd-write-actions")
		b.HackVotingAccounts = viper.GetBool("boot-cmd-hack-voting-accounts")

		if len(args) == 0 {
			b.BootSequenceFile = "boot_sequence.yaml"
		} else {
			b.BootSequenceFile = args[0]
		}

		b.ReuseGenesis = viper.GetBool("boot-cmd-reuse-genesis")

		keyBag, ok := api.Signer.(*eos.KeyBag)
		if ok {
			if len(keyBag.Keys) != 0 {
				key := keyBag.Keys[0]
				b.EphemeralPrivateKey = key
				b.EphemeralPublicKey = key.PublicKey()
			}
		}

		err := b.Boot()
		errorCheck("failed eos-bios boot", err)
	},
}

func init() {
	RootCmd.AddCommand(bootCmd)

	homedir, err := homedir.Dir()
	errorCheck("couldn't find home dir", err)

	bootCmd.Flags().BoolP("reuse-genesis", "", false, "Load genesis data from genesis.json, genesis.pub and genesis.key instead of creating a new one.")
	bootCmd.Flags().StringP("cache-path", "", filepath.Join(homedir, ".eosc-boot-cache"), "directory to store downloaded contract and ABI data")
	bootCmd.Flags().BoolP("hack-voting-accounts", "", false, "This will take accounts with large stakes and put a well known public key in place, so the community can test voting.")
	bootCmd.Flags().BoolP("write-actions", "", false, "Write generated actions to actions.jsonl")
}
