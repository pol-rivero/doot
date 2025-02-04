package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pol-rivero/doot/lib/cache"
	"github.com/pol-rivero/doot/lib/constants"
	"github.com/stretchr/testify/assert"
)

func TestCache_GetInitial(t *testing.T) {
	SetUp(t)
	cacheObj := cache.Load()
	assert.Equal(t, cache.CURRENT_CACHE_VERSION, cacheObj.Version, "Unexpected cache version")
	assert.Empty(t, cacheObj.InstalledDirs, "New cache should have no entries")

	filesCache := cacheObj.UseDir("SomeDir")
	assert.NotNil(t, filesCache, "Unexpected nil files cache")
	assert.Empty(t, filesCache.Targets, "New files cache should have no entries")
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
	assert.FileExists(t, cacheFile(), "Cache file not created")
	bytes, err := os.ReadFile(cacheFile())
	assert.NoError(t, err, "Error reading cache file")
	assert.NotEmpty(t, bytes, "Cache file is empty")

	// Load the cache again and check that the data is the same
	cacheObj = cache.Load()
	filesCache = cacheObj.UseDir("SomeDir")
	assert.ElementsMatch(t, filesCache.Targets, []string{"SomeFile.txt"}, "Unexpected files in cache")
	filesCache = cacheObj.UseDir("AnotherDir")
	assert.Equal(t, []string{"AnotherFile.txt"}, filesCache.Targets, "Unexpected files in cache")
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
	assert.Equal(t, cache.CURRENT_CACHE_VERSION, cacheObj.Version, "Unexpected cache version")
	assert.Empty(t, cacheObj.InstalledDirs, "Expected 0 installed directories")
}

func TestCache_MalformedCache(t *testing.T) {
	SetUp(t)
	err := os.WriteFile(cacheFile(), []byte("This is not a cache file"), 0644)
	assert.NoError(t, err, "Error writing cache file")

	// Load the cache again and check that it was reset
	cacheObj := cache.Load()
	assert.Empty(t, cacheObj.InstalledDirs, "Expected 0 installed directories")
}

func TestCache_DefaultsToHomeCacheDir(t *testing.T) {
	SetUp(t)
	os.Unsetenv(constants.ENV_DOOT_CACHE_DIR)
	cacheObj := cache.Load()
	cacheObj.Save()

	assert.NoFileExists(t, cacheFile(), "Cache unexpectedly saved in unset environment variable")
	assert.FileExists(t, filepath.Join(homeDir(), ".cache", "doot", "doot-cache.bin"), "Cache not saved in default location")
}
