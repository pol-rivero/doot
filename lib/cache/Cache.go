package cache

import (
	"log"
	"os"
)

const CURRENT_CACHE_VERSION uint32 = 1

func FromFile() DootCache {
	fileContents, err := os.ReadFile(getCachePath())
	if err != nil {
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

func (cache *DootCache) SaveToFile() {
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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error retrieving home directory: %v", err)
		os.Exit(1)
	}
	err = os.MkdirAll(homeDir+"/.cache/doot", 0755)
	if err != nil {
		log.Fatalf("Error creating cache directory: %v", err)
		os.Exit(1)
	}
	return homeDir + "/.cache/doot/cache.bin"
}
