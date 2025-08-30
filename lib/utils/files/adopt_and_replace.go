package files

import (
	"os"

	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/log"
	"github.com/pol-rivero/doot/lib/linkmode"
	. "github.com/pol-rivero/doot/lib/types"
)

// https://stackoverflow.com/a/58148921
func ReplaceWithLink(target AbsolutePath, dotfilesSource AbsolutePath, linkMode linkmode.LinkMode) error {
	tempLocation := target.AppendExtension(common.DOOT_BACKUP_EXT)
	if err := os.Remove(tempLocation.Str()); err != nil && !os.IsNotExist(err) {
		log.Error("Failed to remove temporary file %s, consider removing it manually.\n%s", tempLocation, err)
		return err
	}

	if err := linkMode.CreateLink(dotfilesSource, tempLocation); err != nil {
		log.Error("Failed to create link %s -> %s: %s", tempLocation, dotfilesSource, err)
		return err
	}

	if err := os.Rename(tempLocation.Str(), target.Str()); err != nil {
		log.Error("Failed to update %s: %s", target, err)
		os.Remove(tempLocation.Str())
		return err
	}

	return nil
}

func AdoptChanges(target AbsolutePath, dotfilesSource AbsolutePath, linkMode linkmode.LinkMode) error {
	log.Info("Adding changes from %s into %s", target, dotfilesSource)
	if err := CopyFile(target.Str(), dotfilesSource.Str()); err != nil {
		log.Error("Failed to copy file %s to %s: %s", target, dotfilesSource, err)
		return err
	}
	return ReplaceWithLink(target, dotfilesSource, linkMode)
}
