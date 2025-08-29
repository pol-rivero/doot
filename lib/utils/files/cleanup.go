package files

import (
	"os"

	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
)

func RemoveAndCleanup(removeFile AbsolutePath, stopAt AbsolutePath) bool {
	err := os.Remove(removeFile.Str())
	if err == nil {
		CleanupEmptyDir(removeFile.Parent(), stopAt)
		return true
	} else if os.IsNotExist(err) {
		log.Info("Link %s does not exist, it may have been removed manually", removeFile)
	} else {
		log.Error("Failed to remove %s: %s", removeFile, err)
	}
	return false
}

func CleanupEmptyDir(dir AbsolutePath, stopAt AbsolutePath) {
	if dir == stopAt {
		return
	}
	dirEntries, err := os.ReadDir(dir.Str())
	if err != nil {
		log.Warning("Could not clean up %s: %s", dir, err)
		return
	}
	if len(dirEntries) > 0 {
		return
	}
	err = os.Remove(dir.Str())
	if err != nil {
		log.Warning("Could not clean up %s: %s", dir, err)
	} else {
		CleanupEmptyDir(dir.Parent(), stopAt)
	}
}
