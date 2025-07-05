package files

import (
	"os"
	"path/filepath"

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
	targetContent, err := os.ReadFile(target.Str())
	if err != nil {
		log.Error("Failed to read file %s: %s", target, err)
		return err
	}
	if err := os.WriteFile(dotfilesSource.Str(), targetContent, 0644); err != nil {
		log.Error("Failed to write file %s: %s", dotfilesSource, err)
		return err
	}
	return ReplaceWithLink(target, dotfilesSource, linkMode)
}

func EnsureParentDir(target AbsolutePath) bool {
	parentDir := filepath.Dir(target.Str())
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		log.Error("Failed to create directory %s: %s", parentDir, err)
		return false
	}
	return true
}

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
