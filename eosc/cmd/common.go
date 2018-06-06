package cmd

import (
	"fmt"
	"os"

	"github.com/eoscanada/eos-go"
	eosvault "github.com/eoscanada/eosc/vault"
	"github.com/spf13/viper"
)

func setupWallet() (*eosvault.Vault, error) {
	walletFile := viper.GetString("vault-file")
	if _, err := os.Stat(walletFile); err != nil {
		return nil, fmt.Errorf("Wallet file %q missing, ", walletFile)
	}

	vault, err := eosvault.NewVaultFromWalletFile(walletFile)
	if err != nil {
		return nil, fmt.Errorf("loading vault, %s", err)
	}

	boxer, err := eosvault.SecretBoxerForType(vault.SecretBoxWrap)
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

	api := eos.New(viper.GetString("api-url"))

	api.SetSigner(vault.KeyBag)

	return api, nil

}

func api() *eos.API {
	return eos.New(viper.GetString("api-url"))
}

func errorCheck(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
