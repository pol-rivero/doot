package cmd

import (
	"os"

	"github.com/pol-rivero/doot/lib/install"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "doot",
	Short: "A fast and simple dotfiles manager that just gets the job done.\nVersion: " + VERSION_STRING,
	Run: func(cmd *cobra.Command, args []string) {
		SetUpLogger(cmd)
		forceCreate, err := cmd.Flags().GetBool("force-create")
		if err != nil {
			panic(err)
		}
		fullClean, err := cmd.Flags().GetBool("full-clean")
		if err != nil {
			panic(err)
		}
		install.Install(forceCreate, fullClean)
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
	installCmd.Args = cobra.NoArgs
	rootCmd.AddGroup(basicCommandsGroup)
	rootCmd.AddGroup(advancedCommandsGroup)
	rootCmd.AddGroup(otherCommandsGroup)

	rootCmd.Flags().Bool("force-create", false, "Ignore the cache and regenerate all symlinks.")
	rootCmd.Flags().Bool("full-clean", false, "Search and remove all broken symlinks that point to the dotfiles directory, even if they were created by another program. Can be slow.")

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Print additional information to stdout.")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Suppress warnings and errors.")
	rootCmd.SetHelpCommandGroupID(otherCommandsGroup.ID)
	rootCmd.SetCompletionCommandGroupID(otherCommandsGroup.ID)
}
