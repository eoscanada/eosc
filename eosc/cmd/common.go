package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/eoscanada/eos-go"
	eosvault "github.com/eoscanada/eosc/vault"
	"github.com/spf13/viper"
)

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
	vault, err := setupWallet()
	if err != nil {
		fmt.Printf("Error setting up wallet: %s\n", err)
		os.Exit(1)
	}

	api := eos.New(viper.GetString("global-api-url"))

	api.SetSigner(vault.KeyBag)

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
	parts := strings.Split(in, "@")
	if len(parts) > 2 {
		return out, fmt.Errorf("permission %q invalid, use account[@permission]", in)
	}

	if len(parts[0]) > 12 {
		return out, fmt.Errorf("account name %q too long", parts[0])
	}

	out.Actor = eos.AccountName(parts[0])
	out.Permission = eos.PermissionName("active")
	if len(parts) == 2 {
		if len(parts[1]) > 12 {
			return out, fmt.Errorf("permission %q name too long", parts[1])
		}

		out.Permission = eos.PermissionName(parts[1])
	}

	return
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

	fmt.Println("Transaction submitted to the network. Confirm at https://eosquery.com/tx/" + resp.TransactionID)
}
