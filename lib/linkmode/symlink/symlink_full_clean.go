package linkmode_symlink

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/cache"
	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
)

func (l *SymlinkLinkMode) RecalculateCache(dotfilesDir AbsolutePath, scanPath string) []*cache.InstalledFile {
	result := make([]*cache.InstalledFile, 0, 128)
	dotfilesDirInfo := common.StatDotfilesDir(dotfilesDir)
	fullCleanScanRecursive(&result, dotfilesDir, dotfilesDirInfo, scanPath)
	return result
}

func fullCleanScanRecursive(result *[]*cache.InstalledFile, dotfilesDir AbsolutePath, dotfilesDirInfo os.FileInfo, scanPath string) {
	entries, err := os.ReadDir(scanPath)
	if err != nil {
		log.Warning("Skipping '%s' due to error: %v", scanPath, err)
		return
	}
	for _, entry := range entries {
		entryName := entry.Name()
		entryPath := filepath.Join(scanPath, entryName)
		if entry.IsDir() {
			if common.IsDotfilesDir(entry, dotfilesDirInfo) {
				// Symlinks inside the dotfiles dir are part of the repo, not installed links
				continue
			}
			fullCleanScanRecursive(result, dotfilesDir, dotfilesDirInfo, entryPath)
		} else if common.DirEntryIsSymlink(entry) {
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
