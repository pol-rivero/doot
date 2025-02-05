package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pol-rivero/doot/lib/cache"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/stretchr/testify/assert"
)

func assertHomeDirContents(t *testing.T, dir string, expected []string) {
	t.Helper()
	path := filepath.Join(homeDir(), dir)
	fileNames := []string{}
	dirEntries, err := os.ReadDir(path)
	assert.NoError(t, err, "Error reading directory")
	for _, entry := range dirEntries {
		fileNames = append(fileNames, entry.Name())
	}
	assert.ElementsMatch(t, expected, fileNames)
}

func assertHomeLink(t *testing.T, link string, target string) {
	t.Helper()
	linkPath := filepath.Join(homeDir(), link)
	targetPath, err := os.Readlink(linkPath)
	assert.NoError(t, err, "Failed to read link")
	assert.Equal(t, target, targetPath)
}

func assertHomeRegularFile(t *testing.T, path string) {
	t.Helper()
	filePath := filepath.Join(homeDir(), path)
	info, err := os.Lstat(filePath)
	assert.NoError(t, err, "Failed to get file info")
	assert.True(t, info.Mode().IsRegular(), "File is not a regular file")
}

func assertCache(t *testing.T, expectTargets []AbsolutePath) {
	dootCache := cache.Load()
	cacheEntry := dootCache.GetEntry(sourceDir() + ":" + homeDir())
	assert.ElementsMatch(t, cacheEntry.GetTargets(), expectTargets)
}
