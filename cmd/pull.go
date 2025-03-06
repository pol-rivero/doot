package cmd

import (
	"github.com/pol-rivero/doot/lib/pull"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	GroupID: advancedCommandsGroup.ID,
	Use:     "pull",
	Short:   "Pull and install changes from the git remote.",
	Run: func(cmd *cobra.Command, args []string) {
		SetUpLogger(cmd)
		pull.Pull()
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
	pullCmd.Args = cobra.NoArgs
}
