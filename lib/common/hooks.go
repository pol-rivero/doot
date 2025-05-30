package common

import (
	"os"

	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/pol-rivero/doot/lib/utils"
)

func RunHooks(dotfilesDir AbsolutePath, hookName string) {
	hookDir := dotfilesDir.Join(HOOKS_DIR).Join(hookName)
	dirEntries, err := os.ReadDir(hookDir.Str())
	if err != nil {
		log.Info("No hooks found for %s", hookName)
		return
	}
	for _, entry := range dirEntries {
		if entry.IsDir() {
			log.Warning("Unexpected directory (%s) in hooks directory. The hooks directory should only contain files or links to files", entry.Name())
			continue
		}
		hookPath := hookDir.Join(entry.Name())
		err := utils.RunCommand(dotfilesDir, hookPath.Str())
		if err != nil {
			handleError(err, hookName, hookPath)
		}
	}
}

func handleError(err error, hookName string, hookPath AbsolutePath) {
	if os.IsPermission(err) {
		log.Fatal("Permission denied for %s hook. Consider making it executable with 'chmod +x %s'", hookName, hookPath)
	} else {
		log.Fatal("Error running %s hook %s: %v", hookName, hookPath, err)
	}
}
