package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pol-rivero/doot/lib/cache"
	"github.com/pol-rivero/doot/lib/constants"
)

func TestCache_GetInitial(t *testing.T) {
	SetUp(t)
	cacheObj := cache.Load()

	if cacheObj.Version != cache.CURRENT_CACHE_VERSION {
		t.Fatalf("Expected cache version %d, got %d", cache.CURRENT_CACHE_VERSION, cacheObj.Version)
	}
	if len(cacheObj.InstalledDirs) != 0 {
		t.Fatalf("Expected 0 installed directories, got %d", len(cacheObj.InstalledDirs))
	}

	filesCache := cacheObj.UseDir("SomeDir")
	if filesCache == nil {
		t.Fatalf("Unexpected nil files cache")
	}
	if len(filesCache.Targets) != 0 {
		t.Fatalf("Expected 0 files in cache, got %d", len(filesCache.Targets))
	}
}

func TestCache_SaveAndLoad(t *testing.T) {
	SetUp(t)
	cacheObj := cache.Load()
	filesCache := cacheObj.UseDir("SomeDir")
	filesCache.Targets = append(filesCache.Targets, "SomeFile.txt")
	filesCache = cacheObj.UseDir("AnotherDir")
	filesCache.Targets = append(filesCache.Targets, "AnotherFile.txt")
	cacheObj.Save()

	// Check that the cache file was created and is not empty
	if _, err := os.Stat(cacheFile()); err != nil {
		t.Fatalf("Cache file not created")
	}
	bytes, err := os.ReadFile(cacheFile())
	if err != nil {
		t.Fatalf("Error reading cache file: %v", err)
	}
	if len(bytes) == 0 {
		t.Fatalf("Cache file is empty")
	}

	// Load the cache again and check that the data is the same
	cacheObj = cache.Load()
	filesCache = cacheObj.UseDir("SomeDir")
	if len(filesCache.Targets) != 1 {
		t.Fatalf("Expected 1 file in cache, got %d", len(filesCache.Targets))
	}
	if filesCache.Targets[0] != "SomeFile.txt" {
		t.Fatalf("Expected file 'SomeFile.txt' in cache, got '%s'", filesCache.Targets[0])
	}
	filesCache = cacheObj.UseDir("AnotherDir")
	if len(filesCache.Targets) != 1 {
		t.Fatalf("Expected 1 file in cache, got %d", len(filesCache.Targets))
	}
	if filesCache.Targets[0] != "AnotherFile.txt" {
		t.Fatalf("Expected file 'AnotherFile.txt' in cache, got '%s'", filesCache.Targets[0])
	}
}

func TestCache_VersionMismatch(t *testing.T) {
	SetUp(t)
	cacheObj := cache.Load()
	cacheObj.Version = cache.CURRENT_CACHE_VERSION + 1
	filesCache := cacheObj.UseDir("SomeDir")
	filesCache.Targets = append(filesCache.Targets, "SomeFile.txt")

	cacheObj.Save()

	// Load the cache again and check that the version was reset
	cacheObj = cache.Load()
	if cacheObj.Version != cache.CURRENT_CACHE_VERSION {
		t.Fatalf("Expected cache version %d, got %d", cache.CURRENT_CACHE_VERSION, cacheObj.Version)
	}
	if len(cacheObj.InstalledDirs) != 0 {
		t.Fatalf("Expected 0 installed directories, got %d", len(cacheObj.InstalledDirs))
	}
}

func TestCache_MalformedCache(t *testing.T) {
	SetUp(t)
	err := os.WriteFile(cacheFile(), []byte("This is not a cache file"), 0644)
	if err != nil {
		t.Fatalf("Error writing cache file: %v", err)
	}

	// Load the cache again and check that it was reset
	cacheObj := cache.Load()
	if len(cacheObj.InstalledDirs) != 0 {
		t.Fatalf("Expected 0 installed directories, got %d", len(cacheObj.InstalledDirs))
	}
}

func TestCache_DefaultsToHomeCacheDir(t *testing.T) {
	SetUp(t)
	os.Unsetenv(constants.ENV_DOOT_CACHE_DIR)
	cacheObj := cache.Load()
	cacheObj.Save()

	if _, err := os.Stat(cacheFile()); err == nil {
		t.Fatalf("Cache unexpectedly saved in unset environment variable")
	}

	if _, err := os.Stat(filepath.Join(homeDir(), ".cache", "doot", "doot-cache.bin")); err != nil {
		t.Fatalf("Cache not saved in default location")
	}
}
