// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eosc/cli"
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
		vault := mustGetWallet()

		vault.PrintPublicKeys()

		listen(vault)
	},
}

func init() {
	vaultCmd.AddCommand(vaultServeCmd)

	vaultServeCmd.Flags().IntP("port", "", 6666, "Listen port")
	vaultServeCmd.Flags().BoolP("auto-accept", "", false, "Whether to auto accept all signature requests, or to ask for a security code on the command line.")
}

func listen(v *eosvault.Vault) {
	http.HandleFunc("/v1/wallet/get_public_keys", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println("Service /v1/wallet/get_public_keys")

		var out []string
		for _, key := range v.KeyBag.Keys {
			out = append(out, key.PublicKey().String())
		}
		json.NewEncoder(w).Encode(out)
	})

	http.HandleFunc("/v1/wallet/sign_transaction", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Incoming signature request")

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

		fmt.Println("")

		if !viper.GetBool("vault-serve-cmd-auto-accept") {
			res, err := cli.GetConfirmation(`- Enter the code "%d" to allow signature: `)
			if err != nil {
				fmt.Println("sign_transaction: error reading confirmation from command line:", err)
				http.Error(w, "error reading confirmation from command line", 500)
				return
			}

			if !res {
				fmt.Println("sign_transaction: security code invalid, not signing request")
				http.Error(w, "security code invalid, not signing request", 401)
				return
			}
		} else {
			fmt.Println("- Auto-signing request")
		}

		signed, err := v.KeyBag.Sign(tx, chainID, requiredKeys...)
		for _, action := range signed.Transaction.Actions {
			action.SetToServer(false)
		}

		if err != nil {
			http.Error(w, fmt.Sprintf("error signing: %s", err), 500)
			return
		}

		cnt, err := json.Marshal(signed)
		if err != nil {
			http.Error(w, fmt.Sprintf("couldn't marshal output: %s", err), 500)
			return
		}
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(cnt)))
		w.WriteHeader(201)
		_, err = w.Write(cnt)
		if err != nil {
			log.Println("Error writing to socket:", err)
		}
	})

	port := viper.GetInt("vault-serve-cmd-port")
	fmt.Printf("Listening for wallet operations on 127.0.0.1:%d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil); err != nil {
		fmt.Printf("Failed listening on port %d: %s\n", port, err)
	}
}
