package cmd

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"

	eos "github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var toolsNameCodecCmd = &cobra.Command{
	Use:   "name-codec [value]",
	Short: "Convert a value to and from name-encoded strings",
	Long: `EOS name encoding creates strings or up to 12 characters out of uint64 values.

This command allows you to do the conversions.

    eosc tools name-codec 1 --from uint64 --to name
    eosc tools name-codec eosio --from name --to uint64
    eosc tools name-codec eosio --from name --to hex

`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]

		fromCodec := viper.Get("tools-name-codec-cmd-from")
		toCodec := viper.Get("tools-name-codec-cmd-to")

		var value uint64
		switch fromCodec {
		case "hex", "hex_be":
			val, err := hex.DecodeString(input)
			errorCheck("decoding hex", err)
			if len(val) != 8 {
				errorCheck("invalid length", fmt.Errorf("length of hex string expected to be 16 characters, found %d", len(input)))
			}

			switch fromCodec {
			case "hex":
				value = binary.LittleEndian.Uint64(val)
			case "hex_be":
				value = binary.BigEndian.Uint64(val)
			}

		case "name":
			var err error
			value, err = eos.StringToName(input)
			errorCheck("invalid name encoded input", err)

		case "uint64":
			var err error
			value, err = strconv.ParseUint(input, 10, 64)
			errorCheck("invalid uint64 input", err)

		default:
			fmt.Printf("Invalid --from codec %q, possible values: hex, hex_be, uint64, name", fromCodec)
			os.Exit(1)
		}

		cnt := make([]byte, 8)
		switch toCodec {
		case "hex":
			binary.LittleEndian.PutUint64(cnt, value)

			fmt.Println(hex.EncodeToString(cnt))
		case "hex_be":
			binary.BigEndian.PutUint64(cnt, value)

			fmt.Println(hex.EncodeToString(cnt))

		case "name":
			fmt.Println(eos.NameToString(value))

		case "uint64":
			fmt.Println(strconv.FormatUint(value, 10))

		default:
			fmt.Printf("Invalid --to codec %q, possible values: hex, hex_be, uint64, name", fromCodec)
			os.Exit(1)
		}
	},
}

func init() {
	toolsCmd.AddCommand(toolsNameCodecCmd)

	toolsNameCodecCmd.Flags().StringP("from", "", "", "From encoding: name, uint64, hex")
	toolsNameCodecCmd.Flags().StringP("to", "", "", "To encoding: name, uint64, hex")

	for _, flag := range []string{"from", "to"} {
		if err := viper.BindPFlag("tools-name-codec-cmd-"+flag, toolsNameCodecCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}
