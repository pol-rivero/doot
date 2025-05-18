package cmd

import (
	"github.com/pol-rivero/doot/lib/common/log"
	"github.com/spf13/cobra"
)

// Replaced in the build process (release.yaml)
const VERSION = "[[VERSION]]"
const COMMIT = "[[COMMIT]]"
const VERSION_STRING = VERSION + " (" + COMMIT + ")"

var versionCmd = &cobra.Command{
	GroupID: otherCommandsGroup.ID,
	Use:     "version",
	Short:   "Show the version information.",
	Run: func(cmd *cobra.Command, args []string) {
		SetUpLogger(cmd)

		if log.IsVerbose() {
			log.Printlnf("Version: %s", VERSION_STRING)
		} else {
			log.Printlnf("%s", VERSION)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Args = cobra.NoArgs
}
