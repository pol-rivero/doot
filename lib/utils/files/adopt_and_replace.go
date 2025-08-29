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
	targetContent, err := os.ReadFile(target.Str())
	if err != nil {
		log.Error("Failed to read file %s: %s", target, err)
		return err
	}
	// If the dotfile is a symlink, writing to it will replace its target file contents, instead of the symlink itself.
	deleteIfSymlink(dotfilesSource)
	if err := os.WriteFile(dotfilesSource.Str(), targetContent, 0644); err != nil {
		log.Error("Failed to write file %s: %s", dotfilesSource, err)
		return err
	}
	return ReplaceWithLink(target, dotfilesSource, linkMode)
}

func deleteIfSymlink(target AbsolutePath) {
	info, err := os.Lstat(target.Str())
	if err != nil {
		log.Error("Failed to stat %s: %s", target, err)
		return
	}
	if common.IsSymlink(info) {
		if err := os.Remove(target.Str()); err != nil {
			log.Error("Failed to remove symlink %s: %s", target, err)
		} else {
			log.Info("Removed symlink %s so that it can be replaced with the new file", target)
		}
	}
}
