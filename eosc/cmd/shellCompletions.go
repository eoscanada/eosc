package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "shell-completion",
	Short: "Generate shell completions",
}

var genBashCompletion = &cobra.Command{
	Use:   "bash",
	Short: "Generate bash completion file output",
	Run: func(cmd *cobra.Command, args []string) {
		if err := RootCmd.GenBashCompletion(os.Stdout); err != nil {
			log.Fatal(err)
		}
	},
}

var genZshCompletion = &cobra.Command{
	Use:   "zsh",
	Short: "Generate zsh completion file output",
	Run: func(cmd *cobra.Command, args []string) {
		if err := RootCmd.GenZshCompletion(os.Stdout); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(completionCmd)
	completionCmd.AddCommand(genBashCompletion)
	completionCmd.AddCommand(genZshCompletion)
}
