package cmd

import (
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

var systemUnlinkAuthCmd = &cobra.Command{
	Use:   "unlinkauth [your account] [code account] [action name]",
	Short: "Unassign a permission currently active for the given code::action pair",
	Long: `Unassign a permission currently active for the given code::action pair.

This undoes the action of linkauth, please refer to the documentation for linkauth for more details.
`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		account := toAccount(args[0], "your account")
		code := toAccount(args[1], "code account")
		actionName := toActionName(args[2], "action name")

		api := getAPI()
		pushEOSCActions(api, system.NewUnlinkAuth(account, code, actionName))
	},
}

func init() {
	systemCmd.AddCommand(systemUnlinkAuthCmd)
}
