package cli

import (
	"bytes"
	"fmt"
	"strconv"
	"text/tabwriter"

	eos "github.com/eoscanada/eos-go"
	"github.com/ryanuber/columnize"
)

const indentPadding = "      "

func FormatBasicAccountInfo(account *eos.AccountResp, config *columnize.Config) string {
	output := []string{
		fmt.Sprintf("privileged: |%v", account.Privileged),
		fmt.Sprintf("created at: |%v", account.Created),
	}

	if account.LastCodeUpdate.Unix() > 0 {
		output = append(output, fmt.Sprintf("code updated at: |%v", account.LastCodeUpdate))
	}

	return columnize.Format(output, config)
}

func FormatCurrencyStats(stats *eos.GetCurrencyStatsResp, config *columnize.Config) string {
	output := []string{
		fmt.Sprintf("supply: |%v", prettifyAsset(stats.Supply)),
		fmt.Sprintf("max supply: |%v", prettifyAsset(stats.MaxSupply)),
		fmt.Sprintf("issuer: | %v", stats.Issuer),
	}

	return columnize.Format(output, config)
}

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
		fmt.Sprintf("%squota: %s| used: %s",
			indentPadding,
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
		fmt.Sprintf("%sstaked:|%s|(total stake delegated from account to self)",
			indentPadding,
			prettifyAsset(account.SelfDelegatedBandwidth.NetWeight),
		),
		fmt.Sprintf("%sdelegated:|%s|(total stake delegated to account from others)",
			indentPadding,
			prettifyAsset(delegatedNet),
		),
		fmt.Sprintf("%sused:|%s", indentPadding, prettifyBytes(int64(account.NetLimit.Used))),
		fmt.Sprintf("%savailable:|%s", indentPadding, prettifyBytes(int64(account.NetLimit.Available))),
		fmt.Sprintf("%slimit:|%s", indentPadding, prettifyBytes(int64(account.NetLimit.Max))),
	}

	return columnize.Format(output, config)
}

func FormatCPUBandwidth(account *eos.AccountResp, config *columnize.Config) string {
	delegatedCPU := account.TotalResources.CPUWeight.Sub(account.SelfDelegatedBandwidth.CPUWeight)

	output := []string{
		"cpu bandwidth:",
		fmt.Sprintf("%sstaked:|%s|(total stake delegated from account to self)",
			indentPadding,
			prettifyAsset(account.SelfDelegatedBandwidth.CPUWeight),
		),
		fmt.Sprintf("%sdelegated:|%s|(total stake delegated to account from others)",
			indentPadding,
			prettifyAsset(delegatedCPU),
		),
		fmt.Sprintf("%sused:|%s", indentPadding, prettifyTime(int64(account.CPULimit.Used))),
		fmt.Sprintf("%savailable:|%s", indentPadding, prettifyTime(int64(account.CPULimit.Available))),
		fmt.Sprintf("%slimit:|%s", indentPadding, prettifyTime(int64(account.CPULimit.Max))),
	}

	return columnize.Format(output, config)
}

func FormatBalances(account *eos.AccountResp, config *columnize.Config) string {
	if account.CoreLiquidBalance.Symbol.Symbol != "" {
		totalStaked := account.SelfDelegatedBandwidth.NetWeight.Add(account.SelfDelegatedBandwidth.CPUWeight)
		if totalStaked.Symbol != account.CoreLiquidBalance.Symbol {
			totalStaked.Symbol = account.CoreLiquidBalance.Symbol
		}
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
			fmt.Sprintf("%sliquid:|%s", indentPadding, prettifyAsset(account.CoreLiquidBalance)),
			fmt.Sprintf("%sstaked:|%s", indentPadding, prettifyAsset(totalStaked)),
			fmt.Sprintf("%sunstaking:|%s", indentPadding, prettifyAsset(totalUnstaking)),
			fmt.Sprintf("%stotal:|%s", indentPadding, prettifyAsset(total)),
		}

		return columnize.Format(output, config)
	} else {
		return "No liquid balance available"
	}
}

func FormatProducers(account *eos.AccountResp, config *columnize.Config) string {
	accounts := prettifyAccounts(account.VoterInfo.Producers, account)
	output := []string{
		"voted for:",
	}
	output = append(output, accounts...)
	return columnize.Format(output, config)
}

func FormatVoterInfo(account *eos.AccountResp, config *columnize.Config) string {
	proxy := "<none>"
	if len(account.VoterInfo.Proxy) > 0 {
		proxy = string(account.VoterInfo.Proxy)
	}
	output := []string{
		"voter info:",
		fmt.Sprintf("%sproxy:|%s", indentPadding, proxy),
		fmt.Sprintf("%sis proxy:|%v", indentPadding, account.VoterInfo.IsProxy == 1),
		fmt.Sprintf("%sstaked:|%d", indentPadding, account.VoterInfo.Staked),
		fmt.Sprintf("%svote weight:|%f", indentPadding, account.VoterInfo.LastVoteWeight),
		fmt.Sprintf("%sproxied vote weight:|%f", indentPadding, account.VoterInfo.ProxiedVoteWeight),
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
	unit := fmt.Sprintf("%cB", "KMGTPE"[exp])

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
		unit = "h"
	} else if value > 1000000*60 {
		value /= float64(1000000 * 60)
		unit = "m"
	} else if value > 1000000 {
		value /= float64(1000000)
		unit = "s"
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

func prettifyAccounts(accounts []eos.AccountName, account *eos.AccountResp) []string {
	names := []string{}
	if len(accounts) == 0 {
		if len(account.VoterInfo.Proxy) > 0 {
			return []string{fmt.Sprintf("%svotes via proxy: %s", indentPadding, account.VoterInfo.Proxy)}
		}
		return []string{fmt.Sprintf("%s%s", indentPadding, "<not voted>")}
	}
	for _, name := range accounts {
		names = append(names, fmt.Sprintf("%s%s", indentPadding, name))
	}

	return names
}

func rightAlignColumnize(value, unit string) string {
	w := new(tabwriter.Writer)
	bs := bytes.NewBuffer([]byte{})
	// Using tabwriter.Debug to output '|' which is the delimited in columnize
	w.Init(bs, 15, 0, 1, ' ', tabwriter.Debug|tabwriter.AlignRight)
	fmt.Fprintf(w, "%s\t%s", value, unit)
	w.Flush()
	return bs.String()
}
