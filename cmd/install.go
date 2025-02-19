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
		install.Install()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Args = cobra.NoArgs
	installCmd.Flags().Bool("full-clean", false, "Ignore the cache and clean up all broken symlinks that point to the\ndotfiles directory, even if they were created by another program. Can be slow.")
}
