package common

import (
	"os"
	"path/filepath"

	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/pol-rivero/doot/lib/utils"
)

func RunHooks(dotfilesDir AbsolutePath, hookName string) {
	hookDir := filepath.Join(dotfilesDir.Str(), HOOKS_DIR, hookName)
	dirEntries, err := os.ReadDir(hookDir)
	if err != nil {
		log.Info("No hooks found for %s", hookName)
		return
	}
	for _, entry := range dirEntries {
		if entry.IsDir() {
			log.Warning("Unexpected directory (%s) in hooks directory. The hooks directory should only contain files or links to files", entry.Name())
			continue
		}
		hookPath := filepath.Join(hookDir, entry.Name())
		err := utils.RunCommand(dotfilesDir, hookPath)
		if err != nil {
			log.Fatal("Error running hook %s: %v", hookPath, err)
		}
	}
}
