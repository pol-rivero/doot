package cmd

import (
	"github.com/pol-rivero/doot/lib/add"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	GroupID: basicCommandsGroup.ID,
	Use:     "add <file> [file2 ...]",
	Short:   "Move one or more files to the dotfiles directory and symlink them.",
	Run: func(cmd *cobra.Command, args []string) {
		SetUpLogger(cmd)
		isCrypt, err := cmd.Flags().GetBool("crypt")
		if err != nil {
			panic(err)
		}
		hostSpecific, err := cmd.Flags().GetBool("host")
		if err != nil {
			panic(err)
		}
		add.Add(args, isCrypt, hostSpecific)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Args = cobra.MinimumNArgs(1)
	addCmd.ArgAliases = []string{"file"}

	addCmd.Flags().Bool("crypt", false, "Add as a private (encrypted) files.")
	addCmd.Flags().Bool("host", false, "Add as host-specific files. Will be moved to the specific directory\nfor the current host, if it exists.")
}
