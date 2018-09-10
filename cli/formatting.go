package cli

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"

	eos "github.com/eoscanada/eos-go"
	"github.com/ryanuber/columnize"
)

const indentPadding = "      "

func FormatPermissions(account *eos.AccountResp, config *columnize.Config) string {
	output := formatNestedPermission([]string{"permissions:"}, account.Permissions, eos.PermissionName(""), "")
	return columnize.Format(output, config)
}
func formatNestedPermission(in []string, permissions []eos.Permission, showChildsOf eos.PermissionName, indent string) (out []string) {
	out = in
	for _, perm := range permissions {
		if perm.Parent != string(showChildsOf) {
			continue
		}
		permValues := []string{}
		for _, key := range perm.RequiredAuth.Keys {
			permValues = append(permValues, fmt.Sprintf("+%d %s", key.Weight, key.PublicKey))
		}
		for _, acct := range perm.RequiredAuth.Accounts {
			permValues = append(permValues, fmt.Sprintf("+%d %s@%s", acct.Weight, acct.Permission.Actor, acct.Permission.Permission))
		}
		for _, wait := range perm.RequiredAuth.Waits {
			permValues = append(permValues, fmt.Sprintf("+%d wait %d seconds", wait.Weight, wait.WaitSec))
		}
		for i, keyValue := range permValues {
			if i == 0 {
				out = append(out,
					fmt.Sprintf("     %s%q w/%d|:|%s",
						indent,
						perm.PermName,
						perm.RequiredAuth.Threshold,
						keyValue,
					),
				)
			} else {
				out = append(out,
					fmt.Sprintf("     ||%s",
						keyValue,
					),
				)
			}
		}
		out = formatNestedPermission(out, permissions, eos.PermissionName(perm.PermName), indent+indentPadding)
	}
	return out
}

func FormatMemory(account *eos.AccountResp, config *columnize.Config) string {
	output := []string{
		"memory:",
		fmt.Sprintf("     quota: %s| used: %s",
			prettifyBytes(account.RAMQuota),
			prettifyBytes(account.RAMUsage),
		),
	}

	return columnize.Format(output, config)
}

func FormatNetworkBandwidth(account *eos.AccountResp, config *columnize.Config) string {
	delegatedNet := account.TotalResources.NetWeight.Sub(account.SelfDelegatedBandwidth.NetWeight)

	output := []string{
		"net bandwidth:",
		fmt.Sprintf("     staked:|%s|(total stake delegated from account to self)",
			prettifyAsset(account.SelfDelegatedBandwidth.NetWeight),
		),
		fmt.Sprintf("     delegated:|%s|(total stake delegated to account from others)",
			prettifyAsset(delegatedNet),
		),
		fmt.Sprintf("     used:|%s", prettifyBytes(int64(account.NetLimit.Used))),
		fmt.Sprintf("     available:|%s", prettifyBytes(int64(account.NetLimit.Available))),
		fmt.Sprintf("     limit:|%s", prettifyBytes(int64(account.NetLimit.Max))),
	}

	return columnize.Format(output, config)
}

func FormatCPUBandwidth(account *eos.AccountResp, config *columnize.Config) string {
	delegatedCPU := account.TotalResources.CPUWeight.Sub(account.SelfDelegatedBandwidth.CPUWeight)

	output := []string{
		"cpu bandwidth:",
		fmt.Sprintf("     staked:|%s|(total stake delegated from account to self)",
			prettifyAsset(account.SelfDelegatedBandwidth.CPUWeight),
		),
		fmt.Sprintf("     delegated:|%s|(total stake delegated to account from others)",
			prettifyAsset(delegatedCPU),
		),
		fmt.Sprintf("     used:|%s", prettifyTime(int64(account.CPULimit.Used))),
		fmt.Sprintf("     available:|%s", prettifyTime(int64(account.CPULimit.Available))),
		fmt.Sprintf("     limit:|%s", prettifyTime(int64(account.CPULimit.Max))),
	}

	return columnize.Format(output, config)
}

func FormatBalances(account *eos.AccountResp, config *columnize.Config) string {
	totalStaked := account.SelfDelegatedBandwidth.NetWeight.Add(account.SelfDelegatedBandwidth.CPUWeight)
	totalUnstaking := eos.Asset{
		Amount: 0,
		Symbol: account.CoreLiquidBalance.Symbol,
	}
	if account.RefundRequest != nil {
		totalUnstaking = account.RefundRequest.CPUAmount.Add(account.RefundRequest.NetAmount)
	}
	total := totalUnstaking.Add(totalStaked).Add(account.CoreLiquidBalance)

	output := []string{
		fmt.Sprintf("%s balances:", account.CoreLiquidBalance.Symbol.Symbol),
		fmt.Sprintf("     liquid:|%s", prettifyAsset(account.CoreLiquidBalance)),
		fmt.Sprintf("     staked:|%s", prettifyAsset(totalStaked)),
		fmt.Sprintf("     unstaking:|%s", prettifyAsset(totalUnstaking)),
		fmt.Sprintf("     total:|%s", prettifyAsset(total)),
	}

	return columnize.Format(output, config)
}

func FormatProducers(account *eos.AccountResp, config *columnize.Config) string {
	output := []string{
		"producers:",
		fmt.Sprintf("     %s", prettifyAccounts(account.VoterInfo.Producers)),
	}

	return columnize.Format(output, config)
}

func prettifyBytes(b int64) string {
	const u = 1024
	if b < u {
		return rightAlignColumnize(fmt.Sprintf("%d", b), "bytes")
	}
	div, exp := int64(u), 0
	for n := b / u; n >= u; n /= u {
		div *= u
		exp++
	}
	value := float64(b) / float64(div)
	unit := fmt.Sprintf("%ciB", "KMGTPE"[exp])

	precision := 3
	if value >= 100 {
		precision = 1
	} else if value >= 10 {
		precision = 2
	}

	return rightAlignColumnize(strconv.FormatFloat(value, 'f', precision, 64), unit)
}

func prettifyTime(micro int64) string {
	value := float64(micro)
	unit := "μs"
	if value > 1000000*60*60 {
		value /= float64(1000000 * 60 * 60)
		unit = "hr"
	} else if value > 1000000*60 {
		value /= float64(1000000 * 60)
		unit = "min"
	} else if value > 1000000 {
		value /= float64(1000000)
		unit = "sec"
	} else if value > 1000 {
		value /= float64(1000)
		unit = "ms"
	}

	precision := 3
	if value >= 100 {
		precision = 1
	} else if value >= 10 {
		precision = 2
	}

	return rightAlignColumnize(strconv.FormatFloat(value, 'f', precision, 64), unit)
}

func prettifyAsset(w eos.Asset) string {
	const unit = 10000
	formatting := fmt.Sprintf("%%.%df", w.Precision)
	return rightAlignColumnize(fmt.Sprintf(formatting, float64(w.Amount)/float64(unit)), w.Symbol.Symbol)

}

func prettifyAccounts(accounts []eos.AccountName) string {
	names := []string{}
	if len(accounts) == 0 {
		return "<not voted>"
	}
	for _, name := range accounts {
		names = append(names, string(name))
	}

	return strings.Join(names, "|  ")
}

func rightAlignColumnize(value, unit string) string {
	w := new(tabwriter.Writer)
	bs := bytes.NewBuffer([]byte(""))
	// Using tabwriter.Debug to output '|' which is the delimited in columnize
	w.Init(bs, 15, 0, 1, ' ', tabwriter.Debug|tabwriter.AlignRight)
	fmt.Fprintf(w, "%s\t%s", value, unit)
	w.Flush()
	return bs.String()
}
