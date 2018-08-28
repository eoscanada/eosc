package cmd

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	yaml2json "github.com/bronze1man/go-yaml2json"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eos-go/sudo"
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
		return nil, fmt.Errorf("wallet file %q missing: %s", walletFile, err)
	}

	vault, err := eosvault.NewVaultFromWalletFile(walletFile)
	if err != nil {
		return nil, fmt.Errorf("loading vault: %s", err)
	}

	boxer, err := eosvault.SecretBoxerForType(vault.SecretBoxWrap, viper.GetString("global-kms-gcp-keypath"))
	if err != nil {
		return nil, fmt.Errorf("secret boxer: %s", err)
	}

	if err := vault.Open(boxer); err != nil {
		return nil, err
	}

	return vault, nil
}

func attachWallet(api *eos.API) {
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
		errorCheck("specified --permission(s) invalid", err)

		for _, act := range actions {
			act.Authorization = levels
		}
	}

	opts := &eos.TxOptions{}

	if chainID := viper.GetString("global-offline-chain-id"); chainID != "" {
		opts.ChainID = toSHA256Bytes(chainID, "--offline-chain-id")
	}

	if headBlockID := viper.GetString("global-offline-head-block"); headBlockID != "" {
		opts.HeadBlockID = toSHA256Bytes(headBlockID, "--offline-head-block")
	}

	if err := opts.FillFromChain(api); err != nil {
		fmt.Println("Error fetching tapos + chain_id from the chain (specify --offline flags for offline operations):", err)
		os.Exit(1)
	}

	tx := eos.NewTransaction(actions, opts)

	tx = optionallySudoWrap(tx, opts)

	tx.SetExpiration(time.Duration(viper.GetInt("global-expiration")) * time.Second)

	signedTx, packedTx := optionallySignTransaction(tx, opts.ChainID, api)

	optionallyPushTransaction(signedTx, packedTx, opts.ChainID, api)
}

func optionallySudoWrap(tx *eos.Transaction, opts *eos.TxOptions) *eos.Transaction {
	if viper.GetBool("global-sudo-wrap") {
		binTx, err := eos.MarshalBinary(tx)
		errorCheck("binary-packing transaction for sudo wrapping", err)

		return eos.NewTransaction([]*eos.Action{sudo.NewExec(eos.AccountName("eosio"), eos.HexBytes(binTx))}, opts)
	}
	return tx
}

func optionallySignTransaction(tx *eos.Transaction, chainID eos.SHA256Bytes, api *eos.API) (signedTx *eos.SignedTransaction, packedTx *eos.PackedTransaction) {
	if !viper.GetBool("global-skip-sign") {
		textSignKeys := viper.GetStringSlice("global-offline-sign-key")
		if len(textSignKeys) > 0 {
			var signKeys []ecc.PublicKey
			for _, key := range textSignKeys {
				pubKey, err := ecc.NewPublicKey(key)
				errorCheck(fmt.Sprintf("parsing public key %q", key), err)

				signKeys = append(signKeys, pubKey)
			}
			api.SetCustomGetRequiredKeys(func(tx *eos.Transaction) ([]ecc.PublicKey, error) {
				return signKeys, nil
			})
		}

		attachWallet(api)

		var err error
		signedTx, packedTx, err = api.SignTransaction(tx, chainID, eos.CompressionNone)
		errorCheck("signing transaction", err)
	} else {
		signedTx = eos.NewSignedTransaction(tx)
	}
	return signedTx, packedTx
}

func optionallyPushTransaction(signedTx *eos.SignedTransaction, packedTx *eos.PackedTransaction, chainID eos.SHA256Bytes, api *eos.API) {
	outputTrx := viper.GetString("global-output-transaction")

	if outputTrx != "" {
		cnt, err := json.MarshalIndent(signedTx, "", "  ")
		errorCheck("marshalling json", err)

		err = ioutil.WriteFile(outputTrx, cnt, 0644)
		errorCheck("writing output transaction", err)

		for _, act := range signedTx.Actions {
			act.SetToServer(false)
		}

		cnt, err = json.MarshalIndent(signedTx, "", "  ")
		errorCheck("marshalling json", err)

		fmt.Println(string(cnt))
		fmt.Println("---")
		fmt.Printf("Transaction written to %q\n", outputTrx)
		fmt.Printf("Sign offline with: --offline-chain-id=%s\n", hex.EncodeToString(chainID))
		fmt.Println("Above is a pretty-printed representation of the outputted file")
	} else {
		if packedTx == nil {
			fmt.Println("A signed transaction is required if you want to broadcast it. Remove --skip-sign (or add --output-transaction ?)")
			os.Exit(1)
		}

		// TODO: print the traces
		pushTransaction(api, packedTx)
	}
}

func pushTransaction(api *eos.API, packedTx *eos.PackedTransaction) {
	resp, err := api.PushTransaction(packedTx)
	errorCheck("pushing transaction", err)

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

func toName(in, field string) eos.Name {
	name, err := cli.ToName(in)
	if err != nil {
		errorCheck(fmt.Sprintf("invalid name format for %q", field), err)
	}

	return name
}

func toSHA256Bytes(in, field string) eos.SHA256Bytes {
	if len(in) != 64 {
		errorCheck(fmt.Sprintf("%q invalid", field), errors.New("should be 64 hexadecimal characters"))
	}

	bytes, err := hex.DecodeString(in)
	errorCheck(fmt.Sprintf("invalid hex in %q", field), err)

	return bytes
}
