package customcmd

import (
	"os"

	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/pol-rivero/doot/lib/utils"
)

func CustomCommand(command string, args []string) {
	dotfilesDir := common.FindDotfilesDir()
	commandPath := verifyCommand(dotfilesDir, command)
	log.Info("Running custom command '%s' with args: %v", commandPath, args)

	err := utils.RunCommand(dotfilesDir, commandPath.Str(), args...)
	if err != nil {
		handleError(err, commandPath)
	} else {
		log.Info("Custom command '%s' executed successfully", commandPath)
	}
}

func verifyCommand(dotfilesDir AbsolutePath, commandName string) AbsolutePath {
	commandPath := dotfilesDir.Join(common.CUSTOM_COMMANDS_DIR).Join(commandName)
	_, err := os.Stat(commandPath.Str())
	if err == nil {
		return commandPath
	}
	if os.IsNotExist(err) {
		log.Fatal("Command '%s' not recognized. You can define a custom command by creating '%s'", commandName, commandPath)
	}
	log.Fatal("Error reading command file '%s': %v", commandPath, err)
	panic("Unreachable")
}

func handleError(err error, commandPath AbsolutePath) {
	if os.IsPermission(err) {
		log.Fatal("Permission denied for custom command. Consider making it executable with 'chmod +x %s'", commandPath)
	} else {
		log.Fatal("Error running custom command %s: %v", commandPath, err)
	}
}
