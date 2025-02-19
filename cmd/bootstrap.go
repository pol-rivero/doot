package cmd

import (
	"github.com/pol-rivero/doot/lib/bootstrap"
	"github.com/pol-rivero/doot/lib/utils/optional"
	"github.com/spf13/cobra"
)

var bootstrapCmd = &cobra.Command{
	GroupID: advancedCommandsGroup.ID,
	Use:     "bootstrap",
	Short:   "Download and install your dotfiles.",
	Run: func(cmd *cobra.Command, args []string) {
		SetUpLogger(cmd)

		keyFile, err := cmd.Flags().GetString("key")
		if err != nil {
			panic(err)
		}
		bootstrap.Bootstrap(args[0], args[1], optional.WrapString(keyFile))
	},
}

func init() {
	rootCmd.AddCommand(bootstrapCmd)

	bootstrapCmd.Args = cobra.ExactArgs(2)
	bootstrapCmd.ArgAliases = []string{"repo", "dotfiles-dir"}

	bootstrapCmd.Flags().String("key", "", "Path to the private key file to use for decryption.")
}
