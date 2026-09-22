package common

import (
	"io/fs"
	"os"

	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
)

func StatDotfilesDir(dotfilesDir AbsolutePath) os.FileInfo {
	info, err := os.Stat(dotfilesDir.Str())
	if err != nil {
		log.Fatal("Failed to stat dotfiles directory %s: %v", dotfilesDir, err)
	}
	return info
}

// Compares by file identity instead of path, so that it works even if the paths are written differently
func IsDotfilesDir(dirEntry fs.DirEntry, dotfilesDirInfo os.FileInfo) bool {
	info, err := dirEntry.Info()
	return err == nil && os.SameFile(info, dotfilesDirInfo)
}
