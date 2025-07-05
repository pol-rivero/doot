package cmd

import (
	"github.com/pol-rivero/doot/lib/commands/ls"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	GroupID: basicCommandsGroup.ID,
	Use:     "ls",
	Short:   "List the installed (symlinked) dotfiles.",
	Run: func(cmd *cobra.Command, args []string) {
		SetUpLogger(cmd)
		asJson, err := cmd.Flags().GetBool("json")
		if err != nil {
			panic(err)
		}
		ls.ListInstalledFiles(asJson)
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

	lsCmd.Args = cobra.NoArgs
	lsCmd.Flags().Bool("json", false, `Output the result as a JSON object ({"installed_path": "dotfile_path", ...}).`)
}
