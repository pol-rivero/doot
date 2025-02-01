package cache

import (
	"log"
	"os"
	"path"

	"github.com/pol-rivero/doot/lib/constants"
)

const CURRENT_CACHE_VERSION uint32 = 1

func Load() DootCache {
	fileContents, err := os.ReadFile(getCachePath())
	if err != nil {
		log.Printf("Cache read error: %v, creating new cache", err)
		return DootCache{
			Version:       CURRENT_CACHE_VERSION,
			InstalledDirs: []*DotfilesDir{},
		}
	}

	var cacheData DootCache
	err = cacheData.UnmarshalBinary(fileContents)
	if err != nil {
		log.Fatalf("Error parsing cache file: %v", err)
		return DootCache{
			Version:       CURRENT_CACHE_VERSION,
			InstalledDirs: []*DotfilesDir{},
		}
	}

	if cacheData.Version != CURRENT_CACHE_VERSION {
		log.Printf("Cache version mismatch: expected %d, got %d", CURRENT_CACHE_VERSION, cacheData.Version)
		return DootCache{
			Version:       CURRENT_CACHE_VERSION,
			InstalledDirs: []*DotfilesDir{},
		}
	}

	return cacheData
}

func (cache *DootCache) Save() {
	marshalledData, err := cache.MarshalBinary()
	if err != nil {
		log.Fatalf("Error marshalling cache data: %v", err)
		return
	}

	err = os.WriteFile(getCachePath(), marshalledData, 0644)
	if err != nil {
		log.Fatalf("Error saving cache file: %v", err)
	}
}

func (cache *DootCache) UseDir(dotfilesDir string) *InstalledFilesCache {
	for _, installedDir := range cache.InstalledDirs {
		if installedDir.DotfilesPath == dotfilesDir {
			return installedDir.InstalledFiles
		}
	}

	newDir := DotfilesDir{
		DotfilesPath:   dotfilesDir,
		InstalledFiles: &InstalledFilesCache{},
	}
	cache.InstalledDirs = append(cache.InstalledDirs, &newDir)
	return newDir.InstalledFiles
}

func getCachePath() string {
	cacheDir := getCacheContainingDir()
	err := os.MkdirAll(cacheDir, 0755)
	if err != nil {
		log.Fatalf("Error creating cache directory: %v", err)
	}
	return path.Join(cacheDir, "cache.bin")
}

func getCacheContainingDir() string {
	cacheDir := os.Getenv(constants.ENV_DOOT_CACHE_DIR)
	if cacheDir != "" {
		return cacheDir
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error retrieving home directory: %v", err)
		os.Exit(1)
	}
	return path.Join(homeDir, ".cache", "doot")
}
