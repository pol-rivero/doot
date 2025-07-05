package install

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pol-rivero/doot/lib/common/cache"
	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
)

func recalculateCache(installedFilesCache *cache.InstalledFilesCache, dotfilesDir AbsolutePath, scanPath string) {
	installedFilesCache.Links = make([]*cache.InstalledFile, 0, 128)
	fullCleanScanRecursive(&installedFilesCache.Links, dotfilesDir, scanPath)
}

func fullCleanScanRecursive(result *[]*cache.InstalledFile, dotfilesDir AbsolutePath, scanPath string) {
	entries, err := os.ReadDir(scanPath)
	if err != nil {
		log.Warning("Skipping '%s' due to error: %v", scanPath, err)
		return
	}
	for _, entry := range entries {
		entryName := entry.Name()
		entryPath := filepath.Join(scanPath, entryName)
		if entry.IsDir() {
			fullCleanScanRecursive(result, dotfilesDir, entryPath)
		} else if entry.Type()&os.ModeSymlink != 0 {
			target, err := os.Readlink(entryPath)
			if err != nil {
				log.Warning("Failed to read symlink %s: %v", entryPath, err)
				continue
			}
			if !strings.HasPrefix(target, dotfilesDir.Str()) {
				continue
			}
			*result = append(*result, &cache.InstalledFile{
				Path:    entryPath,
				Content: target,
			})
		}
	}
}
