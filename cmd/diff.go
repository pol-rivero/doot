package cmd

import (
	"github.com/pol-rivero/doot/lib/commands/diff"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	GroupID: basicCommandsGroup.ID,
	Use:     "diff",
	Short:   "Show uncommitted changes in the dotfiles repository.",
	Run: func(cmd *cobra.Command, args []string) {
		SetUpLogger(cmd)
		staged, _ := cmd.Flags().GetBool("staged")
		diff.Diff(staged)
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)

	diffCmd.Args = cobra.NoArgs
	diffCmd.Flags().Bool("staged", false, "Show only staged changes.")
}
