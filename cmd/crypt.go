package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cryptCmd = &cobra.Command{
	GroupID: advancedCommandsGroup.ID,
	Use:     "crypt",
	Short:   "Manage private (encrypted) files.",
}

var cryptInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the repository for use with encrypted files.",
	Run: func(cmd *cobra.Command, args []string) {
		SetUpLogger(cmd)
		fmt.Println("crypt init called")
	},
}

var cryptUnlockCmd = &cobra.Command{
	Use:   "unlock [key_file]",
	Short: "Unlock a protected repository to be able to access encrypted files.",
	Run: func(cmd *cobra.Command, args []string) {
		SetUpLogger(cmd)
		fmt.Println("crypt unlock called")
	},
}

var cryptExportKeyCmd = &cobra.Command{
	Use:   "export-key <output_file>",
	Short: "Export the public key to a file.",
	Run: func(cmd *cobra.Command, args []string) {
		SetUpLogger(cmd)
		fmt.Println("crypt export-key called")
	},
}

var cryptAddGpgUserCmd = &cobra.Command{
	Use:   "add-gpg-user <user_id>",
	Short: "Add a GPG user to the list of users that can decrypt the repository.",
	Run: func(cmd *cobra.Command, args []string) {
		SetUpLogger(cmd)
		fmt.Println("crypt add-gpg-user called")
	},
}

func init() {
	rootCmd.AddCommand(cryptCmd)

	cryptCmd.AddCommand(cryptInitCmd)
	cryptCmd.AddCommand(cryptUnlockCmd)
	cryptCmd.AddCommand(cryptExportKeyCmd)
	cryptCmd.AddCommand(cryptAddGpgUserCmd)

	cryptUnlockCmd.Args = cobra.MaximumNArgs(1)
	cryptUnlockCmd.ArgAliases = []string{"key_file"}

	cryptExportKeyCmd.Args = cobra.ExactArgs(1)
	cryptExportKeyCmd.ArgAliases = []string{"output_file"}

	cryptAddGpgUserCmd.Args = cobra.ExactArgs(1)
	cryptAddGpgUserCmd.ArgAliases = []string{"user_id"}
}
