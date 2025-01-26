package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "doot",
	Short: "A fast and simple dotfiles manager that just gets the job done.",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var basicCommandsGroup = &cobra.Group{
	ID:    "basicCommands",
	Title: "Basic commands:",
}

var advancedCommandsGroup = &cobra.Group{
	ID:    "advancedCommands",
	Title: "Advanced commands:",
}

var otherCommandsGroup = &cobra.Group{
	ID:    "otherCommands",
	Title: "Other commands:",
}

func init() {
	rootCmd.AddGroup(basicCommandsGroup)
	rootCmd.AddGroup(advancedCommandsGroup)
	rootCmd.AddGroup(otherCommandsGroup)
	rootCmd.Flags().Bool("full-clean", false, "Ignore the cache and clean up all broken symlinks that point to the\ndotfiles directory, even if they were created by another program. Can be slow.")
	rootCmd.SetHelpCommandGroupID(otherCommandsGroup.ID)
	rootCmd.SetCompletionCommandGroupID(otherCommandsGroup.ID)
}
