package cmd

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	eos "github.com/eoscanada/eos-go"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

var toolsNamesCmd = &cobra.Command{
	Use:   "names value [value ...]",
	Short: "Convert value(s) to and from name-encoded strings",
	Long: `EOS name encoding creates strings or up to 12 characters out of uint64 values.

This command auto-detects encoding and converts it to different encodings.
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		showHeader := len(args) > 1

		if showHeader {
			// Add a starting blank line just like the command with 1 argument
			fmt.Println()
		}

		for _, input := range args {
			if showHeader {
				fmt.Printf("  # %s", input)
			}

			printName(input)
		}
	},
}

func printName(input string) {
	showFrom := map[string]uint64{}

	baseHex, err := hex.DecodeString(input)
	if err == nil {
		if len(baseHex) == 8 {
			showFrom["hex"] = binary.LittleEndian.Uint64(baseHex)
			showFrom["hex_be"] = binary.BigEndian.Uint64(baseHex)
		} else if len(baseHex) == 4 {
			showFrom["hex"] = uint64(binary.LittleEndian.Uint32(baseHex))
			showFrom["hex_be"] = uint64(binary.BigEndian.Uint32(baseHex))
		}
	}

	fromSymbol, err := eos.StringToSymbol(input)
	if err == nil {
		symbolUint, err := fromSymbol.ToUint64()
		if err == nil {
			showFrom["symbol"] = symbolUint
		}
	}

	fromSymbolCode, err := eos.StringToSymbolCode(input)
	if err == nil {
		showFrom["symbol_code"] = uint64(fromSymbolCode)
	}

	fromName, err := eos.StringToName(input)
	if err == nil {
		showFrom["name"] = fromName
	}

	fromUint64, err := strconv.ParseUint(input, 10, 64)
	if err == nil {
		showFrom["uint64"] = fromUint64
	}

	someFound := false
	rows := []string{"| from \\ to | hex | hex_be | hex_rev_u32 | name | uint64 | symbol | symbol_code", "| --------- | --- | ------ | ----------- | ---- | ------ | ------ | ----------- |"}
	for _, from := range []string{"hex", "hex_be", "name", "uint64", "symbol", "symbol_code"} {
		val, found := showFrom[from]
		if !found {
			continue
		}
		someFound = true

		row := []string{from}
		for _, to := range []string{"hex", "hex_be", "hex_rev_u32", "name", "uint64", "symbol", "symbol_code"} {

			cnt := make([]byte, 8)
			switch to {
			case "hex":
				binary.LittleEndian.PutUint64(cnt, val)
				row = append(row, hex.EncodeToString(cnt))
			case "hex_be":
				binary.BigEndian.PutUint64(cnt, val)
				row = append(row, hex.EncodeToString(cnt))

			case "hex_rev_u32":
				if val > math.MaxUint32 {
					row = append(row, "-")
				} else {
					cnt4 := make([]byte, 4)
					binary.BigEndian.PutUint32(cnt4, math.MaxUint32-uint32(val))
					row = append(row, hex.EncodeToString(cnt4))
				}

			case "name":
				row = append(row, eos.NameToString(val))

			case "uint64":
				row = append(row, strconv.FormatUint(val, 10))

			case "symbol":
				row = append(row, symbOrDash(fmt.Sprintf("%d,%s", uint8(val&0xFF), eos.SymbolCode(val>>8).String())))

			case "symbol_code":
				row = append(row, symbOrDash(eos.SymbolCode(val).String()))
			}
		}
		rows = append(rows, "| "+strings.Join(row, " | ")+" |")
	}

	if !someFound {
		fmt.Printf("Couldn't decode %q with any of these methods: hex, hex_be, name, uint64\n", input)
		os.Exit(1)
	}

	fmt.Println("")
	fmt.Println(columnize.SimpleFormat(rows))
	fmt.Println("")
}

var symbOrDashRE = regexp.MustCompile(`^[0-9A-Z,]+$`)

func symbOrDash(input string) string {
	if !symbOrDashRE.MatchString(input) {
		return "-"
	}
	return input
}

func init() {
	toolsCmd.AddCommand(toolsNamesCmd)
}
