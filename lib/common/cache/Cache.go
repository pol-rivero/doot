package cache

import (
	"os"
	"path"

	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
)

const CURRENT_CACHE_VERSION uint32 = 1

func Load() DootCache {
	fileContents, err := os.ReadFile(getCachePath())
	if err != nil {
		log.Info("Cache read error: %v, creating new cache", err)
		return DootCache{
			Version: CURRENT_CACHE_VERSION,
			Entries: []*CacheEntry{},
		}
	}

	var cacheData DootCache
	err = cacheData.UnmarshalBinary(fileContents)
	if err != nil {
		log.Warning("Error parsing cache file: %v, creating new cache", err)
		return DootCache{
			Version: CURRENT_CACHE_VERSION,
			Entries: []*CacheEntry{},
		}
	}

	if cacheData.Version != CURRENT_CACHE_VERSION {
		log.Info("Cache version mismatch: expected %d, got %d", CURRENT_CACHE_VERSION, cacheData.Version)
		return DootCache{
			Version: CURRENT_CACHE_VERSION,
			Entries: []*CacheEntry{},
		}
	}

	return cacheData
}

func (cache *DootCache) Save() {
	marshalledData, err := cache.MarshalBinary()
	if err != nil {
		log.Error("Error marshalling cache data: %v", err)
		return
	}

	err = os.WriteFile(getCachePath(), marshalledData, 0644)
	if err != nil {
		log.Error("Error saving cache file: %v", err)
	}
}

func (cache *DootCache) GetEntry(cacheKey string) *InstalledFilesCache {
	for _, entry := range cache.Entries {
		if entry.CacheKey == cacheKey {
			return entry.InstalledFiles
		}
	}

	newEntry := CacheEntry{
		CacheKey:       cacheKey,
		InstalledFiles: &InstalledFilesCache{},
	}
	cache.Entries = append(cache.Entries, &newEntry)
	return newEntry.InstalledFiles
}

func (filesCache *InstalledFilesCache) GetTargets() []AbsolutePath {
	targets := make([]AbsolutePath, 0, len(filesCache.Targets))
	for _, target := range filesCache.Targets {
		targets = append(targets, NewAbsolutePath(target))
	}
	return targets
}

func (filesCache *InstalledFilesCache) SetTargets(targets []AbsolutePath) {
	filesCache.Targets = make([]string, 0, len(targets))
	for _, target := range targets {
		filesCache.Targets = append(filesCache.Targets, target.Str())
	}
}

func getCachePath() string {
	cacheDir := getCacheContainingDir()
	err := os.MkdirAll(cacheDir, 0755)
	if err != nil {
		log.Fatal("Error creating cache directory: %v", err)
	}
	return path.Join(cacheDir, "doot-cache.bin")
}

func getCacheContainingDir() string {
	cacheDir := os.Getenv(common.ENV_DOOT_CACHE_DIR)
	if cacheDir != "" {
		return cacheDir
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error retrieving home directory: %v", err)
	}
	return path.Join(homeDir, ".cache", "doot")
}
