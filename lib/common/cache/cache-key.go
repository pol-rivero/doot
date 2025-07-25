package cache

import (
	"path/filepath"

	. "github.com/pol-rivero/doot/lib/types"
)

func ComputeCacheKey(dotfilesDir AbsolutePath, targetDir string) string {
	return dotfilesDir.Str() + string(filepath.ListSeparator) + targetDir
}
