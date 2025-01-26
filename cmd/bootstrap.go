package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var bootstrapCmd = &cobra.Command{
	GroupID: advancedCommandsGroup.ID,
	Use:     "bootstrap",
	Short:   "Download and install your dotfiles.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bootstrap called")
	},
}

func init() {
	rootCmd.AddCommand(bootstrapCmd)
	bootstrapCmd.Flags().String("key", "", "Path to the private key file to use for decryption.")
}
