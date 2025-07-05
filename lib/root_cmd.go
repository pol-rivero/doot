package lib

import (
	"github.com/pol-rivero/doot/lib/commands/install"
	"github.com/pol-rivero/doot/lib/customcmd"
	"github.com/spf13/cobra"
)

func ExecuteRootCmd(cmd *cobra.Command, rawArgs []string) {
	isCustomCommand := len(rawArgs) > 0
	if isCustomCommand {
		customcmd.CustomCommand(rawArgs[0], rawArgs[1:])
	} else {
		ExecuteInstall(cmd)
	}
}

func ExecuteInstall(cmd *cobra.Command) {
	fullClean, err := cmd.Flags().GetBool("full-clean")
	if err != nil {
		panic(err)
	}
	install.Install(fullClean)
}
