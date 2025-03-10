package cmd

import (
	"github.com/pol-rivero/doot/lib/install"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	GroupID: basicCommandsGroup.ID,
	Use:     "install",
	Short:   "Install or incrementally update the symlinks. This is the default command.",
	Run: func(cmd *cobra.Command, args []string) {
		SetUpLogger(cmd)
		fullClean, err := cmd.Flags().GetBool("full-clean")
		if err != nil {
			panic(err)
		}
		install.Install(fullClean)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Args = cobra.NoArgs
	installCmd.Flags().Bool("full-clean", false, "Search and remove all broken symlinks that point to the dotfiles directory, even if they were created by another program. Can be slow.")
}
