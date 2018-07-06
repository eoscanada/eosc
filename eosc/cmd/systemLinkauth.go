package cmd

import (
	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var systemLinkAuthCmd = &cobra.Command{
	Use:   "linkauth [your account] [code account] [action name] [permission name]",
	Short: "Assign a permission to the given code::action pair",
	Long: `Assign a permission to the given code::action pair.

By default, accounts have an "owner" and "active" key and with the
"active" key, you can sign all transactions for that account.

By using "updateauth", you can create a new permission with a
different set of keys, account delegation and wait times.
See "eosc system updateauth --help" for details.

Once done, you can use "linkauth" to assign that permission to a
code::action pair. Next time you want to sign a transaction destined
to that code::action, you will need to authorize it with the
associated permission.

This is a way to delegate authority on your account in a granular way,
down to the action level.

EXAMPLE:

In an account with a lots of EOS, you can set a permission called
"accounting" and you give 1 different key to 3 employees in the
accounting department, and set a "waits" of 24h (to for a delay on
transactions, with the option to cancel them if found to be unlawful)

You then set the "eosio.token::transfer" action to be assigned to that
permission.  Now you have delegated the possibility to transfer coins
to the accounting department, but kept all other privileges with the
"active" key (which, if the "accounting" permission has "active" as
parent, still can sign transfers).
`,
	Args: cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		account := toAccount(args[0], "your account")
		code := toAccount(args[1], "code account")
		actionName := eos.ActionName(toName(args[2], "action name"))
		permission := eos.PermissionName(toName(args[3], "permission name"))

		api := getAPI()

		pushEOSCActions(api, system.NewLinkAuth(account, code, actionName, permission))
	},
}

func init() {
	systemCmd.AddCommand(systemLinkAuthCmd)
}
