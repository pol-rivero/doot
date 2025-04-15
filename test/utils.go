package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pol-rivero/doot/lib/common/cache"
	"github.com/pol-rivero/doot/lib/crypt"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/stretchr/testify/assert"
)

func assertHomeDirContents(t *testing.T, dir string, expected []string) {
	t.Helper()
	assertDirContents(t, filepath.Join(homeDir(), dir), expected)
}

func assertSourceDirContents(t *testing.T, dir string, expected []string) {
	t.Helper()
	assertDirContents(t, filepath.Join(sourceDir(), dir), expected)
}

func assertDirContents(t *testing.T, path string, expected []string) {
	t.Helper()
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
	info, err := os.Lstat(linkPath)
	assert.NoError(t, err, "Failed to get link info")
	assert.True(t, info.Mode()&os.ModeSymlink != 0, "File is not a symlink")
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

type AssertCacheEntry struct {
	Path    AbsolutePath
	Content string
}

func assertCache(t *testing.T, expectTargets []AssertCacheEntry) {
	t.Helper()
	dootCache := cache.Load()
	cacheEntry := dootCache.GetEntry(sourceDir() + ":" + homeDir())
	expectMap := make(map[AbsolutePath]AbsolutePath)
	for _, expectTarget := range expectTargets {
		expectMap[AbsolutePath(expectTarget.Path)] = AbsolutePath(expectTarget.Content)
	}
	assertSymlinkCollection(t, cacheEntry.GetLinks(), expectMap)
}

func assertSymlinkCollection(t *testing.T, targets SymlinkCollection, expect map[AbsolutePath]AbsolutePath) {
	t.Helper()
	assert.Len(t, targets.Iter(), len(expect))
	for k, v := range targets.Iter() {
		assert.Contains(t, expect, k)
		assert.Equal(t, v, expect[k])
	}
}

func initializeGitCrypt() {
	createNode(sourceDir(), Dir(".git", []FsNode{
		Dir("git-crypt", []FsNode{
			Dir("keys", []FsNode{
				File("default"),
			}),
		}),
		Dir("info", []FsNode{
			FsFile{
				Name:    "attributes",
				Content: crypt.GetGitAttributesContentForTesting(),
			},
		}),
	}))
}

func createHookFile(hook string, scriptName string, content string) {
	createNode(sourceDir(), Dir("doot", []FsNode{
		Dir("hooks", []FsNode{
			Dir(hook, []FsNode{
				FsFile{
					Name:    scriptName,
					Content: content,
				},
			}),
		}),
	}))
	os.Chmod(filepath.Join(sourceDir(), "doot", "hooks", hook, scriptName), 0755)
}

func createCustomCommandFile(name string, content string) {
	createNode(sourceDir(), Dir("doot", []FsNode{
		Dir("commands", []FsNode{
			FsFile{
				Name:    name,
				Content: content,
			},
		}),
	}))
	os.Chmod(filepath.Join(sourceDir(), "doot", "commands", name), 0755)
}

func readFile(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(content)
}
