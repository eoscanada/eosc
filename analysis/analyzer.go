package analysis

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	eos "github.com/eoscanada/eos-go"

	// Load these so `Unpack` does Action unpacking with known ABIs.
	_ "github.com/eoscanada/eos-go/forum"
	_ "github.com/eoscanada/eos-go/msig"
	"github.com/eoscanada/eos-go/sudo"
	"github.com/eoscanada/eos-go/system"
	_ "github.com/eoscanada/eos-go/token"
)

type Analyzer struct {
	API     *eos.API
	Verbose bool
	Writer  *bytes.Buffer
}

func NewAnalyzer(verbose bool) *Analyzer {
	return &Analyzer{
		Verbose: verbose,
		Writer:  &bytes.Buffer{},
	}
}

func (a *Analyzer) AnalyzePacked(trx *eos.PackedTransaction) (err error) {
	a.Pln()
	a.Pln("---------------------------------------------------------------------")
	a.Pln("------------------------- PACKED TRANSACTION ------------------------")
	a.Pln("---------------------------------------------------------------------")
	a.Pln()
	a.Pf("Transaction ID: %s\n", trx.ID())
	a.Pf("Signatures: %q\n", trx.Signatures)
	a.Pf("Packed context free data length: %d\n", len(trx.PackedContextFreeData))
	a.VerbDump(trx.PackedContextFreeData)
	a.Pf("Packed transaction data length: %d\n", len(trx.PackedTransaction))
	a.VerbDump(trx.PackedContextFreeData)
	a.Pln()
	a.Pln("---------------------------------------------------------------------")
	a.Pln("----------------------- SIGNED TRANSACTION --------------------------")
	a.Pln("---------------------------------------------------------------------")
	a.Pln()

	sTx, err := trx.Unpack()
	if err != nil {
		return
	}

	a.Pf("Number of context-free data blobs (on Transaction): %d\n", len(sTx.ContextFreeData))
	for idx, blob := range sTx.ContextFreeData {
		a.Pf("%d. Blob length: %d\n", idx+1, len(blob))
		a.VerbDump(blob)
	}

	return a.AnalyzeSignedTransaction(sTx)
}

func (a *Analyzer) AnalyzeSignedTransaction(sTx *eos.SignedTransaction) (err error) {
	return a.AnalyzeTransaction(sTx.Transaction)
}
func (a *Analyzer) AnalyzeTransaction(tx *eos.Transaction) (err error) {

	a.Pln()
	a.Pln("---------------------------------------------------------------------")
	a.Pln("----------------------- TRANSACTION HEADER --------------------------")
	a.Pln("---------------------------------------------------------------------")
	a.Pln()

	now := time.Now().UTC()
	a.Pf("Expiration: %s (in %s, analysis time: %s)\n", tx.Expiration.Time, tx.Expiration.Time.Sub(now), now)
	a.Pf("Expiration: %s\n", tx.Expiration.Time)
	a.Pf("Reference block number: %d\n", tx.RefBlockNum)
	a.Pf("Reference block prefix: %x\n", tx.RefBlockPrefix)
	a.Pf("Maximum net usage words (of 8 bytes, 0 = unlimited): %d\n", tx.MaxNetUsageWords)
	a.Pf("Maximum CPU usage in milliseconds (0 = unlimited): %d\n", tx.MaxCPUUsageMS)
	a.Pf("Number of seconds to delay transaction (cancellable during that time): %d\n", tx.DelaySec)

	a.Pln()
	a.Pln("---------------------------------------------------------------------")
	a.Pln("------------------------------ ACTIONS ------------------------------")
	a.Pln("---------------------------------------------------------------------")
	a.Pln()

	a.Pf("Context-free actions: %d\n", len(tx.ContextFreeActions))
	for idx, act := range tx.ContextFreeActions {
		if err := a.analyzeAction(idx, act); err != nil {
			return err
		}
	}

	a.Pln()

	a.Pf("Actions: %d\n", len(tx.Actions))
	for idx, act := range tx.Actions {
		if err := a.analyzeAction(idx, act); err != nil {
			return err
		}
	}

	return nil
}

func (a *Analyzer) analyzeAction(idx int, act *eos.Action) (err error) {
	var auths []string
	for _, auth := range act.Authorization {
		auths = append(auths, fmt.Sprintf("%s@%s", auth.Actor, auth.Permission))
	}
	a.Pf("%d. Action %s::%s, authorized by: %s\n", idx+1, act.Account, act.Name, strings.Join(auths, ", "))

	switch obj := act.ActionData.Data.(type) {
	case *system.SetCode:
		a.Pf("Set code for account: %s\n", obj.Account)
		a.Pf("VM type/version: %d/%d\n", obj.VMType, obj.VMVersion)
		h := sha256.New()
		_, _ = h.Write(obj.Code)
		a.Pf("Code's SHA256: %s\n", hex.EncodeToString(h.Sum(nil)))
		a.Pf("Contains the string 'SYS': %v\n", bytes.Contains(obj.Code, []byte("SYS")))
		a.Pf("Contains the string 'EOS': %v\n", bytes.Contains(obj.Code, []byte("EOS")))
		a.VerbDump(obj.Code)

	case *system.SetABI:
		a.Pf("Set ABI for account: %s\n", obj.Account)
		var unpackedABI eos.ABI
		if err := eos.UnmarshalBinary(obj.ABI, &unpackedABI); err != nil {
			a.Pf("Error: couldn't unpack the ABI therein: %s\n", err)
		}
		jsonABI, err := json.MarshalIndent(unpackedABI, "", "  ")
		if err != nil {
			a.Pf("Error: couldn't serialize ABI into JSON: %s\n", err)
		}
		a.VerbPln("JSON representation of the ABI:")
		a.VerbPf("%s\n", string(jsonABI))

	case *system.SetRAMRate:
		a.Pf("Set increase in RAM to: %d bytes per block\n", obj.BytesPerBlock)

	case *system.SetPriv:
		a.Pf("Set account %q to be privileged=%v\n", obj.Account, obj.IsPriv)

	case *system.NewAccount:
		a.Pf("Create a new account named %q, created by %q with the following authority structure:\n", obj.Name, obj.Creator)
		jsonAuth, err := json.MarshalIndent(map[string]interface{}{
			"owner":  obj.Owner,
			"active": obj.Active,
		}, "", "  ")
		if err != nil {
			a.Pf("Error: couldn't serialize auth structs: %s\n", err)
		}
		a.Pln(string(jsonAuth))

	case *system.DelegateBW:
		a.Pf("Delegate bandwidth from account %q to receiver %q\n", obj.From, obj.Receiver)
		a.Pf("- Network bandwidth stake increase: %s\n", obj.StakeNet)
		a.Pf("- CPU bandwidth stake increase: %s\n", obj.StakeCPU)
		a.Pf("Transfer ownership of stakes? %v\n", obj.Transfer)

	case *system.BuyRAM:
		a.Pf("Account %q is buying RAM for receiver %q, for a total of value of %q\n", obj.Payer, obj.Receiver, obj.Quantity)

	case *system.BuyRAMBytes:
		a.Pf("Account %q is buying RAM for receiver %q, for %d bytes (whatever the market value)\n", obj.Payer, obj.Receiver, obj.Bytes)

	case *system.Setalimits:
		a.Pf("Account %q is to set:\n* RAM limit to %d bytes\n* Net weight: %d bytes\n* CPU weight: %d ms\n", obj.Account, obj.RAMBytes, obj.NetWeight, obj.CPUWeight)

	case *sudo.Exec:

		a.Pf("Account %q executes a transaction impersonating another account\n", obj.Executer)
		var tx *eos.Transaction
		err := eos.UnmarshalBinary(obj.Transaction, &tx)
		if err != nil {
			a.Pf("Error: unpacking sudo transaction: %s\n", err)
		} else {
			a.Pln("-------------------- SUDO TRANSACTION BEGIN -----------------------")
			if err := a.AnalyzeTransaction(tx); err != nil {
				a.Pf("Error: analyzing sudo transaction: %s\n", err)
			}
			a.Pln("-------------------- SUDO TRANSACTION END -----------------------")
		}

	default:
		data, err := a.API.ABIBinToJSON(act.Account, eos.Name(act.Name), act.ActionData.HexData)
		if err != nil {
			a.Pf("Couldn't decode ABI: %s\n", err)
			a.Pf("Hex data: %s\n", hex.EncodeToString(act.ActionData.HexData))
		}
		j, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			a.Pf("Couldn't marshal JSON: %s\n", err)
			a.Pf("Hex data: %s\n", hex.EncodeToString(act.ActionData.HexData))
		}
		a.Pf("Data:\n%s\n", string(j))
		a.Pln()
		return nil
	}
	a.Pln()
	a.Pln()

	return nil
}

// Pln is a short for Println on the Writer
func (a *Analyzer) Pln(v ...interface{}) {
	fmt.Fprintln(a.Writer, v...)
}

// VerbPln is a short for Println on the Writer, in Verbose mode.
func (a *Analyzer) VerbPln(v ...interface{}) {
	if a.Verbose {
		fmt.Fprintln(a.Writer, v...)
	}
}

// VerbDump is a short for spew.Fdump on the Writer, in Verbose mode.
func (a *Analyzer) VerbDump(v ...interface{}) {
	if a.Verbose {
		spew.Fdump(a.Writer, v...)
	}
}

// Dump is a short for spew.Fdump on the Writer.
func (a *Analyzer) Dump(v ...interface{}) {
	spew.Fdump(a.Writer, v...)
}

// Pf is a short for Println on the Writer
func (a *Analyzer) Pf(format string, v ...interface{}) {
	fmt.Fprintf(a.Writer, format, v...)
}

// VerbPf is a short for Println on the Writer, in Verbose mode.
func (a *Analyzer) VerbPf(format string, v ...interface{}) {
	if a.Verbose {
		fmt.Fprintf(a.Writer, format, v...)
	}
}
