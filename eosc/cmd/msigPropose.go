// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/msig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// msigProposeCmd represents the msigPropose command
var msigProposeCmd = &cobra.Command{
	Use:   "propose [proposer] [proposal name] [transaction_file.json]",
	Short: "Propose a new transaction in the eosio.msig contract",
	Long: `Propose a new transaction in the eosio.msig contract

Don't forget '--request' to list authorities needed to approve this
transaction.

The --request-producers option will get the top 30 producers and
require their signatures. This is because there is rotation in the top
21 producers, and 15 out of 21 active producers are required upon
execution (which might differ from the time a proposal was made)

The --with-subaccounts option will navigate the accounts listed in
--request (or --request-producers) and automatically add them to the
requested permission.  This simplifies the multisignature flows.

The --with-owner option, which requires --with-subaccounts, will also
add accounts listed in the owner permissions of the different accounts.

  `,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		proposer := toAccount(args[0], "proposer")
		proposalName := toName(args[1], "proposal name")
		transactionFileName := args[2]

		cnt, err := ioutil.ReadFile(transactionFileName)
		errorCheck("reading transaction file", err)

		var tx *eos.Transaction
		err = json.Unmarshal(cnt, &tx)
		errorCheck("parsing transaction file", err)

		var requested []eos.PermissionLevel
		if viper.GetBool("multisig-propose-cmd-request-producers") {
			out, err := requestProducers(api)
			errorCheck("recursing to get producers accounts", err)

			for el := range out {
				chunks := strings.Split(el, "@")
				requested = append(requested, eos.PermissionLevel{
					Actor:      eos.AccountName(chunks[0]),
					Permission: eos.PermissionName(chunks[1]),
				})
			}

		} else {
			requested, err = permissionsToPermissionLevels(viper.GetStringSlice("multisig-propose-cmd-request"))
			errorCheck("requested permissions", err)
			if len(requested) == 0 {
				errorCheck("--request", errors.New("missing values"))
			}
		}

		if viper.GetBool("multisig-propose-cmd-with-subaccounts") {
			out := make(map[string]bool)
			for _, req := range requested {
				out, err = recurseAccounts(api, out, string(req.Actor), string(req.Permission), 0, 4)
				if err != nil {
					errorCheck("failed recursing", err)
				}
			}

			requested = []eos.PermissionLevel{}
			for el := range out {
				chunks := strings.Split(el, "@")
				requested = append(requested, eos.PermissionLevel{
					Actor:      eos.AccountName(chunks[0]),
					Permission: eos.PermissionName(chunks[1]),
				})
			}
		}

		sort.Slice(requested, func(i, j int) bool {
			el1 := requested[i]
			el2 := requested[j]
			if el1.Actor < el2.Actor {
				return true
			}
			if el1.Actor > el2.Actor {
				return false
			}
			return el1.Permission < el2.Permission
		})

		pushEOSCActions(api,
			msig.NewPropose(proposer, proposalName, requested, tx),
		)
	},
}

func getProducersTable(api *eos.API) (prods producers, err error) {
	lowerBound := ""
	for {
		response, err := api.GetTableRows(
			eos.GetTableRowsRequest{
				Scope:      "eosio",
				Code:       "eosio",
				Table:      "producers",
				JSON:       true,
				LowerBound: lowerBound,
				Limit:      5000,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("get producers table: %s", err)
		}

		var rows producers
		json.Unmarshal(response.Rows, &rows)
		if err != nil {
			return nil, fmt.Errorf("json unmarshal: %s", err)
		}

		prods = append(prods, rows...)

		if !response.More {
			break
		}

		if len(rows) != 0 {
			last := rows[len(rows)-1]
			owner := last["owner"].(string)
			val, _ := eos.StringToName(owner)
			lowerBound = eos.NameToString(val + 1)
		}
	}
	return
}

func requestProducers(api *eos.API) (out map[string]bool, err error) {
	producers, err := getProducersTable(api)
	errorCheck("get producers table", err)

	sort.Slice(producers, producers.Less)

	out = make(map[string]bool)
	for idx, p := range producers {
		if len(out) > 29 {
			break
		}

		newAcct := fmt.Sprintf("%s@active", p["owner"].(string))

		if isActive, _ := p["is_active"].(float64); isActive != 1 {
			fmt.Printf("Skipping inactive no. %d: %s\n", idx+1, newAcct)
			continue
		}

		fmt.Printf("Adding no. %d: %s\n", idx+1, newAcct)
		out[newAcct] = true
	}

	return
}

func recurseAccounts(api *eos.API, in map[string]bool, account string, permission string, level, maxLevels int) (out map[string]bool, err error) {
	out = in

	newAcct := fmt.Sprintf("%s@%s", account, permission)
	if _, found := out[newAcct]; found {
		return
	}

	fmt.Println("      - ADDING:", newAcct)
	out[newAcct] = true

	if level >= maxLevels {
		return out, nil
	}

	//fmt.Println("Fetching account", account)
	resp, err := api.GetAccount(eos.AccountName(account))
	if err != nil {
		return nil, err
	}

	curPerm := permissionByName(resp.Permissions, permission)
	for {
		if !viper.GetBool("multisig-propose-cmd-with-owner") && curPerm.PermName == "owner" {
			break
		}

		if curPerm.PermName == "" {
			break
		}

		for _, acct := range curPerm.RequiredAuth.Accounts {
			out, err = recurseAccounts(api, out, string(acct.Permission.Actor), string(acct.Permission.Permission), level+1, maxLevels)
			if err != nil {
				return nil, err
			}
		}

		curPerm = permissionByName(resp.Permissions, curPerm.Parent)
	}

	return
}

func permissionByName(perms []eos.Permission, name string) eos.Permission {
	for _, perm := range perms {
		if perm.PermName == name {
			return perm
		}
	}
	return eos.Permission{}
}

func init() {
	msigCmd.AddCommand(msigProposeCmd)

	// --request
	// --request-producers
	// --with-subaccounts
	// --with-owner

	msigProposeCmd.Flags().StringSlice("request", []string{}, "Accounts and permissions requested. Comma-separated (or multiple --request) is supported, with 'account' or 'account@permission' (defaults to 'active' permission)")
	msigProposeCmd.Flags().Bool("request-producers", false, "Request permissions from top 30 producers (not just 21 because of potential rotation during long periods of time)")
	msigProposeCmd.Flags().Bool("with-subaccounts", false, "Recursively fetch subaccounts for signature (simplifies your life with multisig accounts)")
	msigProposeCmd.Flags().Bool("with-owner", false, "Also include owner permissions when doing recursion. Requires --with-subaccounts")
}
