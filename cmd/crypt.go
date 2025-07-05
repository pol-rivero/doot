package cmd

import (
	"github.com/pol-rivero/doot/lib/commands/crypt"
	"github.com/pol-rivero/doot/lib/utils/optional"
	"github.com/spf13/cobra"
)

var cryptCmd = &cobra.Command{
	GroupID: advancedCommandsGroup.ID,
	Use:     "crypt",
	Short:   "Manage private (encrypted) files. See wiki for usage.",
}

var cryptInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the repository for use with encrypted files.",
	Run: func(cmd *cobra.Command, args []string) {
		SetUpLogger(cmd)
		crypt.Init()
	},
}

var cryptUnlockCmd = &cobra.Command{
	Use:   "unlock [key_file]",
	Short: "Unlock a protected repository to be able to access encrypted files. If using GPG, the key file is optional. Otherwise, use the key file obtained with 'doot crypt export-key'.",
	Run: func(cmd *cobra.Command, args []string) {
		SetUpLogger(cmd)
		keyFile := optional.Empty[string]()
		if len(args) > 0 {
			keyFile = optional.WrapString(args[0])
		}
		crypt.Unlock(keyFile)
	},
}

var cryptLockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Undo the 'unlock' command and re-encrypt the private files in the work tree.",
	Run: func(cmd *cobra.Command, args []string) {
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			panic(err)
		}
		SetUpLogger(cmd)
		crypt.Lock(force)
	},
}

var cryptExportKeyCmd = &cobra.Command{
	Use:   "export-key <output_file>",
	Short: "Export the private key to a file.",
	Run: func(cmd *cobra.Command, args []string) {
		SetUpLogger(cmd)
		outputFile := args[0]
		crypt.ExportKey(outputFile)
	},
}

var cryptAddGpgUserCmd = &cobra.Command{
	Use:   "add-gpg-user <user_id>",
	Short: "Add the public GPG key to the repository, so that its corresponding private key can be used to unlock the repository.",
	Run: func(cmd *cobra.Command, args []string) {
		SetUpLogger(cmd)
		userId := args[0]
		crypt.AddGpgUser(userId)
	},
}

func init() {
	rootCmd.AddCommand(cryptCmd)

	cryptCmd.AddCommand(cryptInitCmd)
	cryptCmd.AddCommand(cryptUnlockCmd)
	cryptCmd.AddCommand(cryptLockCmd)
	cryptCmd.AddCommand(cryptExportKeyCmd)
	cryptCmd.AddCommand(cryptAddGpgUserCmd)

	cryptInitCmd.Args = cobra.NoArgs

	cryptUnlockCmd.Args = cobra.MaximumNArgs(1)
	cryptUnlockCmd.ArgAliases = []string{"key_file"}

	cryptExportKeyCmd.Args = cobra.ExactArgs(1)
	cryptExportKeyCmd.ArgAliases = []string{"output_file"}

	cryptAddGpgUserCmd.Args = cobra.ExactArgs(1)
	cryptAddGpgUserCmd.ArgAliases = []string{"user_id"}

	cryptLockCmd.Args = cobra.NoArgs
	cryptLockCmd.Flags().Bool("force", false, "Lock even if unclean (you may lose uncommited work)")
}
