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

	passphrase, err := eosvault.GetPassword("Enter passphrase to unlock vault: ")
	if err != nil {
		return nil, fmt.Errorf("reading passphrase: %s", err)

	}

	switch vault.SecretBoxWrap {
	case "passphrase":
		err = vault.OpenWithPassphrase(passphrase)
		if err != nil {
			return nil, fmt.Errorf("reading passphrase: %s", err)
		}
	default:
		return nil, fmt.Errorf("ERROR unsupported secretbox wrapping method: %q", vault.SecretBoxWrap)

	}

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
