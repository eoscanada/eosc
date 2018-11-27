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

Pass --requested-permissions
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
		if viper.GetBool("msig-propose-cmd-request-producers") {
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
			requested, err = permissionsToPermissionLevels(viper.GetStringSlice("msig-propose-cmd-requested-permissions"))
			errorCheck("requested permissions", err)
			if len(requested) == 0 {
				errorCheck("--requested-permissions", errors.New("missing values"))
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

		fmt.Println("MAMA", requested)

		pushEOSCActions(api,
			msig.NewPropose(proposer, proposalName, requested, tx),
		)
	},
}

func requestProducers(api *eos.API) (out map[string]bool, err error) {
	response, err := api.GetTableRows(
		eos.GetTableRowsRequest{
			Scope: "eosio",
			Code:  "eosio",
			Table: "producers",
			JSON:  true,
			Limit: 5000,
		},
	)
	errorCheck("get table rows", err)

	var producers producers
	errorCheck("json unmarshaling producers list", json.Unmarshal(response.Rows, &producers))

	sort.Slice(producers, producers.Less)

	out = make(map[string]bool)
	for idx, p := range producers {
		if idx > 29 {
			break
		}
		fmt.Printf("Recursing producer %d: %s\n", idx+1, p["owner"])
		out, err = recurseAccounts(api, out, p["owner"].(string), "active", 0, 4)
		if err != nil {
			errorCheck("failed recursing", err)
		}
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

	msigProposeCmd.Flags().StringSlice("requested-permissions", []string{}, "Permissions requested, specify multiple times or separated by a comma.")

	msigProposeCmd.Flags().Bool("request-producers", false, "Request permissions from top 30 producers (not just 21 because of potential rotation during long periods of time)")

	for _, flag := range []string{"requested-permissions", "request-producers"} {
		if err := viper.BindPFlag("msig-propose-cmd-"+flag, msigProposeCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}

}
