package utils

import (
	"os"

	"github.com/pol-rivero/doot/lib/constants"
	"github.com/pol-rivero/doot/lib/log"
	. "github.com/pol-rivero/doot/lib/types"
)

func ReplaceWithSymlink(target AbsolutePath, dotfilesSource AbsolutePath) error {
	tempLocation := target.Str() + constants.DOOT_BACKUP_EXT
	err := os.Rename(target.Str(), tempLocation)
	if err != nil {
		log.Error("Failed to move %s to %s: %s", target, tempLocation, err)
		return err
	}

	err = os.Symlink(dotfilesSource.Str(), target.Str())
	if err != nil {
		log.Error("Failed to create link %s -> %s: %s", target, dotfilesSource, err)
		restoreErr := os.Rename(tempLocation, target.Str())
		if restoreErr != nil {
			log.Error("Failed to restore %s from %s! Consider restoring it manually.\n%s", target, tempLocation, restoreErr)
		}
		return err
	}

	err = os.Remove(tempLocation)
	if err != nil {
		log.Warning("Failed to remove temporary file %s, consider removing it manually.\n%s", tempLocation, err)
		return err
	}

	return nil
}
