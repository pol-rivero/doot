package common

import (
	"io/fs"
	"os"
)

func IsSymlink(fileInfo fs.FileInfo) bool {
	return fileInfo.Mode()&os.ModeSymlink != 0
}

func DirEntryIsSymlink(dirEntry fs.DirEntry) bool {
	return dirEntry.Type()&os.ModeSymlink != 0
}
