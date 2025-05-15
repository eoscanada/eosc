package cmd

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eos-go/sudo"
	"github.com/eoscanada/eosc/cli"
	eosvault "github.com/eoscanada/eosc/vault"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const shortFormAuthHelp = `
- An optional threshold for the whole structure: "3=" (defauts to "1=")
- Comma-separated permission levels:
    - a public key
      or:
    - an account name, with optional "@permission" (defaults to "@active")
  For each permission levels, an optional "+2" suffix (defaults to "+1")

EXAMPLES

An authority with a threshold of 1, gated by a single key with a
weight of 1:

    EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV

An authority with a threshold of 1, gated by two accounts each with a weight of 1:

    myaccount,youraccount

An authority with a threshold of 2, gated by two accounts each with a weight of 1

    2=myaccount,youraccount

An authority with a threshold of 3, requiring admin (+2) and one of the two
employees (each +1):

    3=admin+2,employee1,employee2

An authority with a threshold of 3, gated by a key with a weight of 2, an
account with a weight of 3, and another account with a weight of 1:

    3=EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV+2,myaccount@secureperm+3,youraccount
`

func mustGetWallet() *eosvault.Vault {
	vault, err := setupWallet()
	errorCheck("wallet setup", err)
	return vault
}

func setupWallet() (*eosvault.Vault, error) {
	walletFile := viper.GetString("global-vault-file")
	if _, err := os.Stat(walletFile); err != nil {
		return nil, fmt.Errorf("wallet file %q missing: %w", walletFile, err)
	}

	vault, err := eosvault.NewVaultFromWalletFile(walletFile)
	if err != nil {
		return nil, fmt.Errorf("loading vault: %w", err)
	}

	boxer, err := eosvault.SecretBoxerForType(vault.SecretBoxWrap, viper.GetString("global-kms-gcp-keypath"))
	if err != nil {
		return nil, fmt.Errorf("secret boxer: %w", err)
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
	httpHeaders := viper.GetStringSlice("global-http-header")
	api := eos.New(sanitizeAPIURL(viper.GetString("global-api-url")))

	for i := 0; i < 25; i++ {
		if val := os.Getenv(fmt.Sprintf("EOSC_GLOBAL_HTTP_HEADER_%d", i)); val != "" {
			httpHeaders = append(httpHeaders, val)
		}
	}

	for _, header := range httpHeaders {
		headerArray := strings.SplitN(header, ": ", 2)
		if len(headerArray) != 2 || strings.Contains(headerArray[0], " ") {
			errorCheck("validating http headers", fmt.Errorf("invalid HTTP Header format"))
		}
		api.Header.Add(headerArray[0], headerArray[1])
	}

	if viper.GetBool("global-allow-partial-signature") {
		api.UsePartialRequiredKeys()
	}
	api.Debug = viper.GetBool("global-debug")

	return api
}

var (
	coreSymbolIsCached bool
	coreSymbol         eos.Symbol
)

func getCoreSymbol() eos.Symbol {
	if coreSymbolIsCached {
		return coreSymbol
	}

	// In the event of a failure, we do not want to re-perform an API call,
	// so let's record the fact that getCoreSymbol is cached right here.
	// The init core symbol will take care of setting an approriate core
	// symbol from global flag and reporting the error.
	coreSymbolIsCached = true
	if err := initCoreSymbol(); err != nil {
		coreSymbol = eos.EOSSymbol
		zlog.Debug(
			"unable to retrieve core symbol from API, falling back to default",
			zap.Error(err),
			zap.Stringer("default", coreSymbol),
		)
	}

	return coreSymbol
}

func initCoreSymbol() error {
	resp, err := getAPI().GetTableRows(context.Background(), eos.GetTableRowsRequest{
		Code:  "eosio",
		Scope: "eosio",
		Table: "rammarket",
		JSON:  true,
	})
	if err != nil {
		return fmt.Errorf("unable to fetch table: %w", err)
	}

	result := gjson.GetBytes(resp.Rows, "0.quote.balance")
	if !result.Exists() {
		return errors.New("table has not expected format")
	}

	asset, err := eos.NewAssetFromString(result.String())
	if !result.Exists() {
		return fmt.Errorf("quote balance asset %q is not valid: %w", result.String(), err)
	}

	zlog.Debug("Retrieved core symbol from API, using it as default core symbol", zap.Stringer("symbol", asset.Symbol))
	coreSymbol = asset.Symbol
	return nil
}

func sanitizeAPIURL(input string) string {
	return strings.TrimRight(input, "/")
}

func errorCheck(prefix string, err error) {
	if err != nil {
		fmt.Printf("ERROR: %s: %s\n", prefix, err)
		if strings.HasSuffix(err.Error(), "connection refused") && strings.Contains(err.Error(), defaultAPIURL) {
			fmt.Println("Have you selected a valid EOS HTTP endpoint ? You can use the --api-url flag or EOSC_GLOBAL_API_URL environment variable.")
		}
		os.Exit(1)
	}
}

func exitWitMessage(message string) {
	fmt.Printf("ERROR: %s\n", message)
	os.Exit(1)
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

func pushEOSCActions(ctx context.Context, api *eos.API, actions ...*eos.Action) {
	pushEOSCActionsAndContextFreeActions(ctx, api, nil, actions)
}

func pushEOSCActionsAndContextFreeActions(ctx context.Context, api *eos.API, contextFreeActions []*eos.Action, actions []*eos.Action) {
	for _, act := range contextFreeActions {
		act.Authorization = nil
	}

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

	if delaySec := viper.GetInt("global-delay-sec"); delaySec != 0 {
		opts.DelaySecs = uint32(delaySec)
	}

	if err := opts.FillFromChain(ctx, api); err != nil {
		fmt.Println("Error fetching tapos + chain_id from the chain (specify --offline flags for offline operations):", err)
		os.Exit(1)
	}

	tx := eos.NewTransaction(actions, opts)
	if len(contextFreeActions) > 0 {
		tx.ContextFreeActions = contextFreeActions
	}

	tx = optionallySudoWrap(tx, opts)

	signedTx, packedTx := optionallySignTransaction(ctx, tx, opts.ChainID, api, true)

	optionallyPushTransaction(ctx, signedTx, packedTx, opts.ChainID, api)
}

func optionallySudoWrap(tx *eos.Transaction, opts *eos.TxOptions) *eos.Transaction {
	if viper.GetBool("global-sudo-wrap") {
		return eos.NewTransaction([]*eos.Action{sudo.NewExec(eos.AccountName("eosio.wrap"), *tx)}, opts)
	}
	return tx
}

func optionallySignTransaction(ctx context.Context, tx *eos.Transaction, chainID eos.SHA256Bytes, api *eos.API, resetExpiration bool) (signedTx *eos.SignedTransaction, packedTx *eos.PackedTransaction) {
	if !viper.GetBool("global-skip-sign") {
		textSignKeys := viper.GetStringSlice("global-offline-sign-key")
		if len(textSignKeys) > 0 {
			var signKeys []ecc.PublicKey
			for _, key := range textSignKeys {
				pubKey, err := ecc.NewPublicKey(key)
				errorCheck(fmt.Sprintf("parsing public key %q", key), err)

				signKeys = append(signKeys, pubKey)
			}
			api.SetCustomGetRequiredKeys(func(ctx context.Context, tx *eos.Transaction) ([]ecc.PublicKey, error) {
				return signKeys, nil
			})
		}

		attachWallet(api)

		if resetExpiration {
			tx.SetExpiration(time.Duration(viper.GetInt("global-expiration")) * time.Second)
		}

		var err error
		signedTx, packedTx, err = api.SignTransaction(ctx, tx, chainID, eos.CompressionNone)
		errorCheck("signing transaction", err)
	} else {
		tx.SetExpiration(time.Duration(viper.GetInt("global-expiration")) * time.Second)

		signedTx = eos.NewSignedTransaction(tx)
	}
	return signedTx, packedTx
}

func optionallyPushTransaction(ctx context.Context, signedTx *eos.SignedTransaction, packedTx *eos.PackedTransaction, chainID eos.SHA256Bytes, api *eos.API) {
	writeTrx := viper.GetString("global-write-transaction")

	if writeTrx != "" {
		cnt, err := json.MarshalIndent(signedTx, "", "  ")
		errorCheck("marshalling json", err)

		annotatedCnt, err := sjson.Set(string(cnt), "chain_id", hex.EncodeToString(chainID))
		errorCheck("adding chain_id", err)

		if writeTrx != "-" {
			err = ioutil.WriteFile(writeTrx, []byte(annotatedCnt), 0o644)
			errorCheck("writing output transaction", err)

			fmt.Printf("Transaction written to %q\n", writeTrx)
		} else {
			os.Stdin.Write([]byte(annotatedCnt))
		}
	} else {
		if packedTx == nil {
			fmt.Println("A signed transaction is required if you want to broadcast it. Remove --skip-sign (or add --write-transaction ?)")
			os.Exit(1)
		}

		// TODO: print the traces
		pushTransaction(ctx, api, packedTx, chainID)
	}
}

func pushTransaction(ctx context.Context, api *eos.API, packedTx *eos.PackedTransaction, chainID eos.SHA256Bytes) {
	resp, err := api.PushTransaction(ctx, packedTx)
	if err != nil {
		if typedErr, ok := err.(eos.APIError); ok {
			printError(typedErr)
		}
		errorCheck("pushing transaction", err)
	}

	trxURL := transactionURL(chainID, resp.TransactionID)
	fmt.Printf("\nTransaction submitted to the network.\n  %s\n", trxURL)
	if resp.BlockID != "" {
		blockURL := blockURL(chainID, resp.BlockID)
		fmt.Printf("Server says transaction was included in block %d:\n  %s\n", resp.Processed.BlockNum, blockURL)
	}
}

func transactionURL(chainID eos.SHA256Bytes, trxID string) string {
	hexChain := hex.EncodeToString(chainID)
	switch hexChain {
	case "aca376f206b8fc25a6ed44dbdc66547c36c6c33e3a119ffbeaef943642f0e906":
		return fmt.Sprintf("https://eosq.app/tx/%s", trxID)
	case "5fff1dae8dc8e2fc4d5b23b2c7665c97f9e9d8edf2b6485a86ba311c25639191":
		return fmt.Sprintf("https://kylin.eosq.app/tx/%s", trxID)
	case "e70aaab8997e1dfce58fbfac80cbbb8fecec7b99cf982a9444273cbc64c41473":
		return fmt.Sprintf("https://jungle.eosq.app/tx/%s", trxID)
	}

	return trxID
}

func blockURL(chainID eos.SHA256Bytes, blockID string) string {
	hexChain := hex.EncodeToString(chainID)
	switch hexChain {
	case "aca376f206b8fc25a6ed44dbdc66547c36c6c33e3a119ffbeaef943642f0e906":
		return fmt.Sprintf("https://eosq.app/block/%s", blockID)
	case "5fff1dae8dc8e2fc4d5b23b2c7665c97f9e9d8edf2b6485a86ba311c25639191":
		return fmt.Sprintf("https://kylin.eosq.app/block/%s", blockID)
	case "e70aaab8997e1dfce58fbfac80cbbb8fecec7b99cf982a9444273cbc64c41473":
		return fmt.Sprintf("https://jungle.eosq.app/block/%s", blockID)
	}

	return blockID
}

func printError(err eos.APIError) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(err)
}

func toAccount(in, field string) eos.AccountName {
	acct, err := cli.ToAccountName(in)
	errorCheck(fmt.Sprintf("invalid account format for %q", field), err)

	return acct
}

func toAsset(symbol eos.Symbol, in, field string) eos.Asset {
	asset, err := eos.NewFixedSymbolAssetFromString(symbol, in)
	errorCheck(fmt.Sprintf("invalid %q value %q", field, in), err)

	return asset
}

func toAssetWithDefaultCoreSymbol(in, field string) eos.Asset {
	if len(strings.Split(in, " ")) == 1 {
		return toCoreAsset(in, field)
	}

	asset, err := eos.NewAssetFromString(in)
	errorCheck(fmt.Sprintf("invalid asset value %q for %q", in, field), err)

	return asset
}

func toCoreAsset(in, field string) eos.Asset {
	return toAsset(getCoreSymbol(), in, field)
}

func toEOSAsset(in, field string) eos.Asset {
	return toAsset(eos.EOSSymbol, in, field)
}

func toREXAsset(in, field string) eos.Asset {
	return toAsset(eos.REXSymbol, in, field)
}

func toName(in, field string) eos.Name {
	name, err := cli.ToName(in)
	if err != nil {
		errorCheck(fmt.Sprintf("invalid name format for %q", field), err)
	}

	return name
}

func toPermissionLevel(in, field string) eos.PermissionLevel {
	perm, err := permissionToPermissionLevel(in)
	if err != nil {
		errorCheck(fmt.Sprintf("invalid permission level for %q", field), err)
	}
	return perm
}

func toActionName(in, field string) eos.ActionName {
	return eos.ActionName(toName(in, field))
}

func toUint16(in, field string) uint16 {
	value, err := strconv.ParseUint(in, 10, 16)
	errorCheck(fmt.Sprintf("invalid uint16 number for %q", field), err)

	return uint16(value)
}

func toUint64(in, field string) uint64 {
	value, err := strconv.ParseUint(in, 10, 64)
	errorCheck(fmt.Sprintf("invalid uint64 number for %q", field), err)

	return value
}

func toSHA256Bytes(in, field string) eos.SHA256Bytes {
	if len(in) != 64 {
		errorCheck(fmt.Sprintf("%q invalid", field), errors.New("should be 64 hexadecimal characters"))
	}

	bytes, err := hex.DecodeString(in)
	errorCheck(fmt.Sprintf("invalid hex in %q", field), err)

	return bytes
}

func isStubABI(abi eos.ABI) bool {
	return abi.Version == "" &&
		abi.Actions == nil &&
		abi.ErrorMessages == nil &&
		abi.Extensions == nil &&
		abi.RicardianClauses == nil &&
		abi.Structs == nil && abi.Tables == nil &&
		abi.Types == nil
}
