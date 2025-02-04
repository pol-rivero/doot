package cmd

import (
	"github.com/pol-rivero/doot/lib/install"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	GroupID: basicCommandsGroup.ID,
	Use:     "clean",
	Short:   "Remove all symlinks created by doot.",
	Run: func(cmd *cobra.Command, args []string) {
		SetUpLogger(cmd)
		install.Clean()
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().Bool("full-clean", false, "Ignore the cache and remove all symlinks that point to the dotfiles directory,\neven if they were created by another program. Can be slow.")
}
