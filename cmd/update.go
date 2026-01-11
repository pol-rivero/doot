package cmd

import (
	"github.com/pol-rivero/doot/lib/commands/update"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	GroupID: basicCommandsGroup.ID,
	Use:     "update",
	Short:   "Interactively review and stage changes in the dotfiles repository.",
	Long: `Interactively review each changed file in your dotfiles repository.
For each file, you can:
  y - Stage this change
  s - Skip this file  
  a - Stage all remaining files
  q - Quit review`,
	Run: func(cmd *cobra.Command, args []string) {
		SetUpLogger(cmd)
		update.Update()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Args = cobra.NoArgs
}
