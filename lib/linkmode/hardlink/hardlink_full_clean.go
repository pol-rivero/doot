package linkmode_hardlink

import (
	"os"
	"path/filepath"

	"github.com/pol-rivero/doot/lib/common/cache"
	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
)

type HardlinkInfo struct {
	Inode uint64
	Dev   uint64
}

type HardlinkInfoMap map[HardlinkInfo]AbsolutePath

func (l *HardlinkLinkMode) RecalculateCache(dotfilesDir AbsolutePath, scanPath string) []*cache.InstalledFile {
	result := make([]*cache.InstalledFile, 0, 128)
	dotfilesDirHardlinks := make(HardlinkInfoMap, 128)
	getHardlinksInfoRecursive(&dotfilesDirHardlinks, dotfilesDir)
	fullCleanScanRecursive(&result, &dotfilesDirHardlinks, scanPath)
	return result
}

func fullCleanScanRecursive(result *[]*cache.InstalledFile, dotfilesDirHardlinks *HardlinkInfoMap, scanPath string) {
	entries, err := os.ReadDir(scanPath)
	if err != nil {
		log.Warning("Skipping '%s' due to error: %v", scanPath, err)
		return
	}
	for _, entry := range entries {
		entryName := entry.Name()
		entryPath := filepath.Join(scanPath, entryName)
		if entry.IsDir() {
			fullCleanScanRecursive(result, dotfilesDirHardlinks, entryPath)
		} else {
			hardlinkInfo := getHardlinkInfo(entryPath)
			if hardlinkInfo == nil {
				continue
			}
			dotfilePath, found := (*dotfilesDirHardlinks)[*hardlinkInfo]
			if !found {
				continue
			}
			*result = append(*result, &cache.InstalledFile{
				Path:    entryPath,
				Content: dotfilePath.Str(),
			})
		}
	}
}

func getHardlinksInfoRecursive(dotfilesDirHardlinks *HardlinkInfoMap, dotfilesDir AbsolutePath) {
	entries, err := os.ReadDir(dotfilesDir.Str())
	if err != nil {
		log.Warning("Skipping '%s' due to error: %v", dotfilesDir, err)
		return
	}
	for _, entry := range entries {
		entryName := entry.Name()
		entryPath := dotfilesDir.Join(entryName)
		if entry.IsDir() {
			getHardlinksInfoRecursive(dotfilesDirHardlinks, entryPath)
		} else {
			hardlinkInfo := getHardlinkInfo(entryPath.Str())
			if hardlinkInfo != nil {
				(*dotfilesDirHardlinks)[*hardlinkInfo] = entryPath
			}
		}
	}
}

func getHardlinkInfo(path string) *HardlinkInfo {
	info, err := unixStat(path)
	if err != nil {
		log.Info("Failed to get hardlink info for %s: %v", path, err)
		return nil
	}
	if info.Nlink <= 1 {
		return nil
	}
	return &HardlinkInfo{
		Inode: info.Ino,
		Dev:   info.Dev,
	}
}
