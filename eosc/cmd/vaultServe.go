// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	eosvault "github.com/eoscanada/eosc/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// vaultServeCmd represents the serve command
var vaultServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serves signing queries on a local port.",
	Long: `Serve will start listening on a local port, offering a
keosd-compatible interface, ready to sign transactions.

It is to be used with tools such as 'cleos' or 'eos-vote' that need
transactions signed before submitting them to an EOS network.
`,
	Run: func(cmd *cobra.Command, args []string) {
		vault, err := setupWallet()
		if err != nil {
			fmt.Println("ERROR:", err)
			os.Exit(1)
		}

		vault.PrintPublicKeys()

		listen(vault)
	},
}

func init() {
	vaultCmd.AddCommand(vaultServeCmd)

	vaultServeCmd.Flags().IntP("port", "p", 6666, "Listen port")

	for _, flag := range []string{"port"} {
		if err := viper.BindPFlag(flag, vaultServeCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}

func listen(v *eosvault.Vault) {
	http.HandleFunc("/v1/wallet/get_public_keys", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Handling get_public_keys")

		var out []string
		for _, key := range v.KeyBag.Keys {
			out = append(out, key.PublicKey().String())
		}
		json.NewEncoder(w).Encode(out)
	})

	http.HandleFunc("/v1/wallet/sign_transaction", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Handling sign_transaction")

		var inputs []json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&inputs); err != nil {
			fmt.Println("sign_transaction: error:", err)
			http.Error(w, "couldn't decode input", 500)
			return
		}

		var tx *eos.SignedTransaction
		var requiredKeys []ecc.PublicKey
		var chainID eos.HexBytes

		if len(inputs) != 3 {
			http.Error(w, "invalid length of message, should be 3 parameters", 500)
			return
		}

		err := json.Unmarshal(inputs[0], &tx)
		if err != nil {
			http.Error(w, "decoding transaction", 500)
			return
		}

		err = json.Unmarshal(inputs[1], &requiredKeys)
		if err != nil {
			http.Error(w, "decoding required keys", 500)
			return
		}

		err = json.Unmarshal(inputs[2], &chainID)
		if err != nil {
			http.Error(w, "decoding chain id", 500)
			return
		}

		signed, err := v.KeyBag.Sign(tx, chainID, requiredKeys...)
		for _, action := range signed.Transaction.Actions {
			action.SetToServer(false)
		}

		if err != nil {
			http.Error(w, fmt.Sprintf("error signing: %s", err), 500)
			return
		}

		w.WriteHeader(201)
		_ = json.NewEncoder(w).Encode(signed)
	})

	// when trying to import a key, we should instruct them to use the `eosc vault` command
	// to create and import new keys.
	//
	// TODO: maybe have a `--writable`  option, so we keep the passphrase, or the parts
	// or the shamir secret, and then write back to the .json file..
	//
	// http.HandleFunc("/v1/wallet/import_key", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("Handling import_key")
	// 	var inputs []string
	// 	_ = json.NewDecoder(r.Body).Decode(&inputs)
	// 	// We're ignoring inputs[0] which is the name of the wallet ("default" by default)

	// 	keyBag.Add(inputs[1])

	// 	w.WriteHeader(201)
	// 	w.Write([]byte("{}"))
	// })

	port := viper.GetInt("port")
	fmt.Printf("Listening for wallet operations on 127.0.0.1:%d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil); err != nil {
		log.Println("Failed listening on port %d: %s\n", port, err)
	}
}
