package common

import (
	"io/fs"
	"os"

	. "github.com/pol-rivero/doot/lib/types"
)

func IsSymlink(fileInfo fs.FileInfo) bool {
	return fileInfo.Mode()&os.ModeSymlink != 0
}

func DirEntryIsSymlink(dirEntry fs.DirEntry) bool {
	return dirEntry.Type()&os.ModeSymlink != 0
}

func IsSymlinkWithTarget(possiblySymlinkPath AbsolutePath, expectedTarget string) bool {
	linkSource, err := os.Readlink(possiblySymlinkPath.Str())
	return err == nil && linkSource == expectedTarget
}
