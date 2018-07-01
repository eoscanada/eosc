package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	yaml2json "github.com/bronze1man/go-yaml2json"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eosc/cli"
	eosvault "github.com/eoscanada/eosc/vault"
	"github.com/spf13/viper"
)

func mustGetWallet() *eosvault.Vault {
	vault, err := setupWallet()
	errorCheck("wallet setup", err)
	return vault
}

func setupWallet() (*eosvault.Vault, error) {
	walletFile := viper.GetString("global-vault-file")
	if _, err := os.Stat(walletFile); err != nil {
		return nil, fmt.Errorf("Wallet file %q missing, ", walletFile)
	}

	vault, err := eosvault.NewVaultFromWalletFile(walletFile)
	if err != nil {
		return nil, fmt.Errorf("loading vault, %s", err)
	}

	boxer, err := eosvault.SecretBoxerForType(vault.SecretBoxWrap, viper.GetString("global-kms-gcp-keypath"))
	if err != nil {
		return nil, fmt.Errorf("secret boxer, %s", err)
	}

	if err := vault.Open(boxer); err != nil {
		return nil, err
	}

	return vault, nil
}

func apiWithWallet() *eos.API {
	api := getAPI()

	walletURLs := viper.GetStringSlice("global-wallet-url")
	if len(walletURLs) == 0 {
		vault, err := setupWallet()
		errorCheck("setting up wallet", err)

		api.SetSigner(vault.KeyBag)
	} else {
		if len(walletURLs) == 1 {
			// If a `walletURLs` has a Username in the path, use instead of `default`.
			api.SetSigner(eos.NewWalletSigner(eos.New(walletURLs[0]), "default"))
		} else {
			fmt.Println("Multi-signer not yet implemented.  Please choose only one `--wallet-url`")
			os.Exit(1)
		}
	}

	return api

}

func getAPI() *eos.API {
	return eos.New(viper.GetString("global-api-url"))
}

func errorCheck(prefix string, err error) {
	if err != nil {
		fmt.Printf("ERROR: %s: %s\n", prefix, err)
		os.Exit(1)
	}
}

func permissionToPermissionLevel(in string) (out eos.PermissionLevel, err error) {
	return eos.NewPermissionLevel(in)
}

func permissionsToPermissionLevels(in []string) (out []eos.PermissionLevel, err error) {
	// loop all parameters
	for _, singleArg := range in {

		// if they specified "account@active,account2", handle that too..
		for _, val := range strings.Split(singleArg, ",") {
			level, err := permissionToPermissionLevel(strings.TrimSpace(val))
			if err != nil {
				return out, err
			}

			out = append(out, level)
		}
	}

	return
}

func pushEOSCActions(api *eos.API, actions ...*eos.Action) {
	permissions := viper.GetStringSlice("global-permission")
	if len(permissions) != 0 {
		levels, err := permissionsToPermissionLevels(permissions)
		if err != nil {
			fmt.Println("Specified permissions invalid:", err)
			os.Exit(1)
		}

		for _, act := range actions {
			act.Authorization = levels
		}
	}

	resp, err := api.SignPushActions(actions...)
	if err != nil {
		fmt.Println("Error signing/pushing transaction:", err)
		os.Exit(1)
	}

	// TODO: print the traces

	//fmt.Println("Transaction submitted to the network. Confirm at https://eosquery.com/tx/" + resp.TransactionID)
	fmt.Println("Transaction submitted to the network. Transaction ID: " + resp.TransactionID)
}

func yamlUnmarshal(cnt []byte, v interface{}) error {
	jsonCnt, err := yaml2json.Convert(cnt)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonCnt, v)
}

func loadYAMLOrJSONFile(filename string, v interface{}) error {
	cnt, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if strings.HasSuffix(strings.ToLower(filename), ".json") {
		return json.Unmarshal(cnt, v)
	}
	return yamlUnmarshal(cnt, v)
}

func toAccount(in, field string) eos.AccountName {
	acct, err := cli.ToAccountName(in)
	if err != nil {
		errorCheck(fmt.Sprintf("invalid account format for %q", field), err)
	}

	return acct
}
