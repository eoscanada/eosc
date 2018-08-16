package cmd

import (
	"fmt"
	"os"
	"time"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/msig"
	"github.com/eoscanada/eos-go/system"
	"github.com/eoscanada/eos-go/token"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var toolsSellAccountCmd = &cobra.Command{
	Use:   "sell-account [sold account] [buyer account] [beneficiary account] [amount]",
	Short: "Create a multisig transaction that both parties need to approve in order to do an atomic sale of your account.",
	Long: `Create a multisig transaction that both parties need to approve in order to do an atomic sale of your account.

Transfers both "owner" and "active" authority to a clone of the buyer's account's authority.

MAKE SURE TO INSPECT THE GENERATED MULTISIG TRANSACTION BEFORE APPROVING IT.
`,
	Args: cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {

		soldAccount := toAccount(args[0], "sold account")
		buyerAccount := toAccount(args[1], "buyer account")
		beneficiaryAccount := toAccount(args[2], "beneficiary account")
		saleAmount, err := eos.NewEOSAssetFromString(args[3])
		errorCheck(`sale "amount" invalid`, err)
		proposalName := viper.GetString("tools-sell-account-cmd-proposal-name")
		memo := viper.GetString("tools-sell-account-cmd-memo")

		api := getAPI()

		soldAccountData, err := api.GetAccount(soldAccount)
		errorCheck("could not find sold account on chain: "+string(soldAccount), err)

		if len(soldAccountData.Permissions) > 2 {
			fmt.Println("WARNING: your account has more than 2 permissions.")
			fmt.Println("This operation hands off control of `owner` and `active` keys.")
			fmt.Println("Please clean-up your permissions before selling your account.")
			os.Exit(1)
		}

		buyerAccountData, err := api.GetAccount(buyerAccount)
		errorCheck("could not find buyer's account on chain", err)

		_, err = api.GetAccount(beneficiaryAccount)
		errorCheck("could not find beneficiary's account on chain", err)

		buyerPermText := viper.GetString("tools-sell-account-cmd-buyer-permission")
		if buyerPermText == "" {
			buyerPermText = string(buyerAccount)
		}
		buyerPerm, err := eos.NewPermissionLevel(buyerPermText)
		errorCheck(`invalid "buyer-permission"`, err)

		myPermText := viper.GetString("tools-sell-account-cmd-seller-permission")
		if myPermText == "" {
			myPermText = string(soldAccount)
		}
		myPerm, err := eos.NewPermissionLevel(myPermText)
		errorCheck(`invalid "seller-permission"`, err)

		targetOwnerAuth, err := sellAccountFindAuthority(buyerAccountData, "owner")
		errorCheck("error finding buyer's owner permission", err)
		targetActiveAuth, err := sellAccountFindAuthority(buyerAccountData, "active")
		errorCheck("error finding buyer's owner permission", err)

		infoResp, err := api.GetInfo()
		errorCheck("couldn't get_info from chain", err)

		tx := eos.NewTransaction([]*eos.Action{
			system.NewUpdateAuth(soldAccount, eos.PermissionName("owner"), eos.PermissionName(""), targetOwnerAuth, eos.PermissionName("owner")),
			system.NewUpdateAuth(soldAccount, eos.PermissionName("active"), eos.PermissionName("owner"), targetActiveAuth, eos.PermissionName("active")),
			token.NewTransfer(buyerAccount, beneficiaryAccount, saleAmount, memo),
		}, &eos.TxOptions{HeadBlockID: infoResp.HeadBlockID})
		tx.SetExpiration(viper.GetDuration("tools-sell-account-cmd-sale-expiration"))

		fmt.Println("Submitting `eosio.msig` proposal:")
		fmt.Printf("  proposer: %s\n", soldAccount)
		fmt.Printf("  proposal_name: %s\n", proposalName)
		fmt.Println("If this transaction is successful, have the other party approve and execute the multisig proposal to an atomic swap.")
		fmt.Println("Review this proposal with:")
		fmt.Printf("  eosc multisig review %s %s", soldAccount, proposalName)
		fmt.Println("")
		msigPermissions := []eos.PermissionLevel{buyerPerm, myPerm, eos.PermissionLevel{Actor: soldAccount, Permission: eos.PermissionName("owner")}}
		pushEOSCActions(api, msig.NewPropose(soldAccount, eos.Name(proposalName), msigPermissions, tx))

	},
}

func init() {
	toolsCmd.AddCommand(toolsSellAccountCmd)

	toolsSellAccountCmd.Flags().StringP("memo", "", "", "Memo message to attach to transfer")
	toolsSellAccountCmd.Flags().StringP("proposal-name", "", "sellaccount", "Proposal name to use in the eosio.msig contract")
	toolsSellAccountCmd.Flags().StringP("buyer-permission", "", "", "Permission required of the buyer (to authorized 'eosio.token::transfer')")
	toolsSellAccountCmd.Flags().StringP("seller-permission", "", "", "Permission required of the seller (you, to authorize 'eosio::updateauth')")
	toolsSellAccountCmd.Flags().DurationP("sale-expiration", "", 1*time.Hour, "Expire proposed transaction after this amount of time (30m, 1h, etc..)")

	for _, flag := range []string{"memo", "seller-permission", "buyer-permission", "proposal-name", "sale-expiration"} {
		if err := viper.BindPFlag("tools-sell-account-cmd-"+flag, toolsSellAccountCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}

func sellAccountFindAuthority(data *eos.AccountResp, targetPerm string) (eos.Authority, error) {
	for _, perm := range data.Permissions {
		if perm.PermName == targetPerm {
			return perm.RequiredAuth, nil
		}
	}
	return eos.Authority{}, fmt.Errorf("permission %q not found in account %q", targetPerm, data.AccountName)
}
