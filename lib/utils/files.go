package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pol-rivero/doot/lib/constants"
	"github.com/pol-rivero/doot/lib/log"
	. "github.com/pol-rivero/doot/lib/types"
)

// https://stackoverflow.com/a/58148921
func ReplaceWithSymlink(target AbsolutePath, dotfilesSource AbsolutePath) error {
	tempLocation := target.Str() + constants.DOOT_BACKUP_EXT
	if err := os.Remove(tempLocation); err != nil && !os.IsNotExist(err) {
		log.Error("Failed to remove temporary file %s, consider removing it manually.\n%s", tempLocation, err)
		return err
	}

	if err := os.Symlink(dotfilesSource.Str(), tempLocation); err != nil {
		log.Error("Failed to create link %s -> %s: %s", tempLocation, dotfilesSource, err)
		return err
	}

	if err := os.Rename(tempLocation, target.Str()); err != nil {
		log.Error("Failed to update %s: %s", target, err)
		os.Remove(tempLocation)
		return err
	}

	return nil
}

func GetTopLevelDir(filePath RelativePath) string {
	filePathStr := string(filePath)
	firstSeparatorIndex := strings.IndexRune(filePathStr, filepath.Separator)
	if firstSeparatorIndex == -1 {
		return filePathStr
	}
	return filePathStr[:firstSeparatorIndex]
}
