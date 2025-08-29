package test

import (
	"os"
	"testing"

	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/cache"
	"github.com/stretchr/testify/assert"
)

func TestCache_GetInitial(t *testing.T) {
	SetUp(t, true)
	cacheObj := cache.Load()
	assert.Equal(t, cache.CURRENT_CACHE_VERSION, cacheObj.Version, "Unexpected cache version")
	assert.Empty(t, cacheObj.Entries, "New cache should have no entries")

	filesCache := cacheObj.GetEntry("cacheKey1")
	assert.NotNil(t, filesCache, "Unexpected nil files cache")
	assert.Empty(t, filesCache.Links, "New files cache should have no entries")
}

func TestCache_SaveAndLoad(t *testing.T) {
	SetUp(t, true)
	cacheObj := cache.Load()
	filesCache := cacheObj.GetEntry("cacheKey1")
	filesCache.Links = append(filesCache.Links, &cache.InstalledFile{
		Path:    "SomeFile.txt",
		Content: "SomeContent",
	})
	filesCache = cacheObj.GetEntry("cacheKey2")
	filesCache.Links = append(filesCache.Links, &cache.InstalledFile{
		Path:    "AnotherFile.txt",
		Content: "AnotherContent",
	})
	cacheObj.Save()

	// Check that the cache file was created and is not empty
	assert.FileExists(t, cacheFile(), "Cache file not created")
	bytes, err := os.ReadFile(cacheFile())
	assert.NoError(t, err, "Error reading cache file")
	assert.NotEmpty(t, bytes, "Cache file is empty")

	// Load the cache again and check that the data is the same
	cacheObj = cache.Load()
	filesCache = cacheObj.GetEntry("cacheKey1")
	assert.ElementsMatch(t, filesCache.Links, []*cache.InstalledFile{
		{
			Path:    "SomeFile.txt",
			Content: "SomeContent",
		},
	})
	filesCache = cacheObj.GetEntry("cacheKey2")
	assert.ElementsMatch(t, filesCache.Links, []*cache.InstalledFile{
		{
			Path:    "AnotherFile.txt",
			Content: "AnotherContent",
		},
	})
	filesCache = cacheObj.GetEntry("wrongKey")
	assert.Empty(t, filesCache.Links, "Expected empty targets for wrong key")

}

func TestCache_VersionMismatch(t *testing.T) {
	SetUp(t, true)
	cacheObj := cache.Load()
	cacheObj.Version = cache.CURRENT_CACHE_VERSION + 1
	filesCache := cacheObj.GetEntry("cacheKey1")
	filesCache.Links = append(filesCache.Links, &cache.InstalledFile{
		Path:    "SomeFile.txt",
		Content: "SomeContent",
	})

	cacheObj.Save()

	// Load the cache again and check that the version was reset
	cacheObj = cache.Load()
	assert.Equal(t, cache.CURRENT_CACHE_VERSION, cacheObj.Version, "Unexpected cache version")
	assert.Empty(t, cacheObj.Entries)
}

func TestCache_MalformedCache(t *testing.T) {
	SetUp(t, true)
	err := os.WriteFile(cacheFile(), []byte("This is not a cache file"), 0644)
	assert.NoError(t, err, "Error writing cache file")

	// Load the cache again and check that it was reset
	cacheObj := cache.Load()
	assert.Empty(t, cacheObj.Entries)
}

func TestCache_DefaultsToHomeCacheDir(t *testing.T) {
	SetUp(t, true)
	os.Unsetenv(common.ENV_DOOT_CACHE_DIR)
	cacheObj := cache.Load()
	cacheObj.Save()

	assert.NoFileExists(t, cacheFile(), "Cache unexpectedly saved in unset environment variable")
	assert.FileExists(t, homeDir()+"/.cache/doot/doot-cache.bin", "Cache not saved in default location")
}
