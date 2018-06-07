package cmd

import (
	"fmt"
	"os"

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

	vault.Open(boxer)

	return vault, nil
}

func apiWithWallet() (*eos.API, error) {
	vault, err := setupWallet()
	if err != nil {
		return nil, err
	}

	api := eos.New(viper.GetString("global-api-url"))

	api.SetSigner(vault.KeyBag)

	return api, nil

}

func api() *eos.API {
	return eos.New(viper.GetString("global-api-url"))
}

func errorCheck(prefix string, err error) {
	if err != nil {
		fmt.Printf("ERROR: %s: %s\n", prefix, err)
		os.Exit(1)
	}
}
