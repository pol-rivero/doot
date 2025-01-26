package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	GroupID: basicCommandsGroup.ID,
	Use:     "add <file> [file2 ...]",
	Short:   "Move one or more files to the dotfiles directory and symlink them.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Args = cobra.MinimumNArgs(1)
	addCmd.ArgAliases = []string{"file"}

	addCmd.Flags().Bool("crypt", false, "Add as a private (encrypted) files.")
	addCmd.Flags().Bool("host", false, "Add as host-specific files. Will be moved to the specific directory\nfor the current host, if it exists.")
}
