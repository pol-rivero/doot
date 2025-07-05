package cmd

import (
	"github.com/pol-rivero/doot/lib/commands/install"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	GroupID: basicCommandsGroup.ID,
	Use:     "clean",
	Short:   "Remove all symlinks created by doot.",
	Run: func(cmd *cobra.Command, args []string) {
		SetUpLogger(cmd)
		fullClean, err := cmd.Flags().GetBool("full-clean")
		if err != nil {
			panic(err)
		}
		install.Clean(fullClean)
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	cleanCmd.Args = cobra.NoArgs
	cleanCmd.Flags().Bool("full-clean", false, "Search and remove all broken symlinks that point to the dotfiles directory, even if they were created by another program. Can be slow.")
}
