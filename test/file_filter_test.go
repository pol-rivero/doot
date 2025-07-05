package test

import (
	"testing"

	"github.com/pol-rivero/doot/lib/commands/install"
	"github.com/pol-rivero/doot/lib/common/config"
	"github.com/pol-rivero/doot/lib/common/glob_collection"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/stretchr/testify/assert"
)

func TestFileFilter_CreateFilter1(t *testing.T) {
	config := &config.Config{
		ExcludeFiles: []string{"**/.*", "file1"},
		IncludeFiles: []string{"file2"},
	}
	filter := install.CreateFilter(config, false)
	assert.True(t, filter.IgnoreHidden, "Expected IgnoreHidden to be true")
	assert.False(t, filter.IgnoreDootCrypt, "Expected IgnoreDootCrypt to be false")
	assert.Equal(t, filter.ExcludeGlobs.Len(), 1, "Expected 1 exclude glob")
	assert.Equal(t, filter.IncludeGlobs.Len(), 1, "Expected 1 include glob")
}

func TestFileFilter_CreateFilter2(t *testing.T) {
	invalid_glob := "*["
	config := &config.Config{
		ExcludeFiles: []string{"file1", "*.txt"},
		IncludeFiles: []string{"file2", invalid_glob},
	}
	filter := install.CreateFilter(config, true)
	assert.False(t, filter.IgnoreHidden, "Expected IgnoreHidden to be false")
	assert.True(t, filter.IgnoreDootCrypt, "Expected IgnoreDootCrypt to be true")
	assert.Equal(t, filter.ExcludeGlobs.Len(), 2, "Expected 2 exclude globs")
	assert.Equal(t, filter.IncludeGlobs.Len(), 1, "Expected 1 include glob (the other one is invalid)")
}

var SCAN_DIR_TESTS = []func(*testing.T){
	func(t *testing.T) {
		filter := install.FileFilter{
			IgnoreHidden:        false,
			IgnoreDootCrypt:     false,
			ExploreExcludedDirs: false,
			ExcludeGlobs:        glob_collection.NewGlobCollection([]string{}),
			IncludeGlobs:        glob_collection.NewGlobCollection([]string{}),
		}
		files := install.ScanDirectory(sourceDirPath(), &filter)
		expectedFiles := []RelativePath{
			"file1",
			"file2",
			"dir1/nestedDir/file3",
			"dir1/nestedDir/.nestedHiddenFile1",
			".hiddenFile",
			".hiddenDir/file4",
			".hiddenDir/.nestedHiddenFile2",
			".hiddenDir/dïrWìthÜnicóde/155helloö/fileABC",
			"secret1.doot-crypt.txt",
			"secret2.doot-crypt",
			"secret-dir1.doot-crypt/file5",
			"secret-dir2.doot-crypt.d/nested.doot-crypt/file6",
			"secret-dir2.doot-crypt.d/nested.doot-crypt/file7.doot-crypt",
			"secret-dir2.doot-crypt.d/nested.doot-crypt/.hiddenInSecretDir",
		}
		assert.ElementsMatch(t, expectedFiles, files, "Unexpected files")
	},

	func(t *testing.T) {
		filter := install.FileFilter{
			IgnoreHidden:        true,
			IgnoreDootCrypt:     false,
			ExploreExcludedDirs: false,
			ExcludeGlobs:        glob_collection.NewGlobCollection([]string{}),
			IncludeGlobs:        glob_collection.NewGlobCollection([]string{}),
		}
		files := install.ScanDirectory(sourceDirPath(), &filter)
		expectedFiles := []RelativePath{
			"file1",
			"file2",
			"dir1/nestedDir/file3",
			"secret1.doot-crypt.txt",
			"secret2.doot-crypt",
			"secret-dir1.doot-crypt/file5",
			"secret-dir2.doot-crypt.d/nested.doot-crypt/file6",
			"secret-dir2.doot-crypt.d/nested.doot-crypt/file7.doot-crypt",
		}
		assert.ElementsMatch(t, expectedFiles, files, "Unexpected files")
	},

	func(t *testing.T) {
		filter := install.FileFilter{
			IgnoreHidden:        false,
			IgnoreDootCrypt:     true,
			ExploreExcludedDirs: false,
			ExcludeGlobs:        glob_collection.NewGlobCollection([]string{}),
			IncludeGlobs:        glob_collection.NewGlobCollection([]string{}),
		}
		files := install.ScanDirectory(sourceDirPath(), &filter)
		expectedFiles := []RelativePath{
			"file1",
			"file2",
			"dir1/nestedDir/file3",
			"dir1/nestedDir/.nestedHiddenFile1",
			".hiddenFile",
			".hiddenDir/file4",
			".hiddenDir/.nestedHiddenFile2",
			".hiddenDir/dïrWìthÜnicóde/155helloö/fileABC",
		}
		assert.ElementsMatch(t, expectedFiles, files, "Unexpected files")
	},

	func(t *testing.T) {
		filter := install.FileFilter{
			IgnoreHidden:        true,
			IgnoreDootCrypt:     true,
			ExploreExcludedDirs: false,
			ExcludeGlobs:        glob_collection.NewGlobCollection([]string{}),
			IncludeGlobs:        glob_collection.NewGlobCollection([]string{}),
		}
		files := install.ScanDirectory(sourceDirPath(), &filter)
		expectedFiles := []RelativePath{
			"file1",
			"file2",
			"dir1/nestedDir/file3",
		}
		assert.ElementsMatch(t, expectedFiles, files, "Unexpected files")
	},

	func(t *testing.T) {
		filter := install.FileFilter{
			IgnoreHidden:        false,
			IgnoreDootCrypt:     false,
			ExploreExcludedDirs: false,
			ExcludeGlobs:        glob_collection.NewGlobCollection([]string{"secret*"}),
			IncludeGlobs:        glob_collection.NewGlobCollection([]string{"*.txt", "**/file6"}),
		}
		files := install.ScanDirectory(sourceDirPath(), &filter)
		assert.NotContains(t, files, RelativePath("secret2.doot-crypt"), "Excluded because it starts with 'secret'")
		assert.Contains(t, files, RelativePath("secret1.doot-crypt.txt"), "Included because ends with .txt")
		assert.NotContains(t, files, RelativePath("secret-dir2.doot-crypt.d/nested.doot-crypt/file6"), "Even though it's included, excluded directories are not explored")
	},

	func(t *testing.T) {
		filter := install.FileFilter{
			IgnoreHidden:        false,
			IgnoreDootCrypt:     false,
			ExploreExcludedDirs: false,
			ExcludeGlobs:        glob_collection.NewGlobCollection([]string{"secret*/**"}),
			IncludeGlobs:        glob_collection.NewGlobCollection([]string{"secret*/nested.doot-crypt", "**/file6"}),
		}
		files := install.ScanDirectory(sourceDirPath(), &filter)
		assert.NotContains(t, files, RelativePath("secret-dir1.doot-crypt/file5"), "Excluded because it starts with 'secret'")
		assert.Contains(t, files, RelativePath("secret-dir2.doot-crypt.d/nested.doot-crypt/file6"), "Now file6 should be returned, because the directory contents are excluded, not the directory itself")
		assert.NotContains(t, files, RelativePath("secret-dir2.doot-crypt.d/nested.doot-crypt/file7.doot-crypt"), "File7 is not included")
	},

	func(t *testing.T) {
		filter := install.FileFilter{
			IgnoreHidden:        false,
			IgnoreDootCrypt:     false,
			ExploreExcludedDirs: false,
			ExcludeGlobs:        glob_collection.NewGlobCollection([]string{"secret*/**"}),
			IncludeGlobs:        glob_collection.NewGlobCollection([]string{"secret*/nested.doot-crypt**"}),
		}
		files := install.ScanDirectory(sourceDirPath(), &filter)
		assert.NotContains(t, files, RelativePath("secret-dir1.doot-crypt/file5"), "Excluded because it starts with 'secret'")
		assert.Contains(t, files, RelativePath("secret-dir2.doot-crypt.d/nested.doot-crypt/file6"), "All children of nested.doot-crypt should be included")
		assert.Contains(t, files, RelativePath("secret-dir2.doot-crypt.d/nested.doot-crypt/file7.doot-crypt"), "All children of nested.doot-crypt should be included")
	},

	func(t *testing.T) {
		filter := install.FileFilter{
			IgnoreHidden:        false,
			IgnoreDootCrypt:     false,
			ExploreExcludedDirs: true,
			ExcludeGlobs:        glob_collection.NewGlobCollection([]string{"secret*", "dir1"}),
			IncludeGlobs:        glob_collection.NewGlobCollection([]string{"**/file6", "dir1/nestedDir"}),
		}
		files := install.ScanDirectory(sourceDirPath(), &filter)
		assert.Contains(t, files, RelativePath("secret-dir2.doot-crypt.d/nested.doot-crypt/file6"), "Excluded directories should be explored")
		assert.Contains(t, files, RelativePath("dir1/nestedDir/file3"), "Excluded directories should be explored")
	},

	func(t *testing.T) {
		filter := install.FileFilter{
			IgnoreHidden:        false,
			IgnoreDootCrypt:     false,
			ExploreExcludedDirs: false,
			ExcludeGlobs:        glob_collection.NewGlobCollection([]string{"**/file2", ".hiddenDir/**/file4", "dir1/nestedDir/**"}),
			IncludeGlobs:        glob_collection.NewGlobCollection([]string{"dir1/nestedDir/**/file3"}),
		}
		files := install.ScanDirectory(sourceDirPath(), &filter)
		assert.Contains(t, files, RelativePath("file1"), "Should not have excluded file1")
		assert.NotContains(t, files, RelativePath("file2"), "Should have excluded file2 (** should also match depth 0)")
		assert.NotContains(t, files, RelativePath(".hiddenDir/file4"), "Should have excluded this file (** should also match depth 0)")
		assert.Contains(t, files, RelativePath(".hiddenDir/.nestedHiddenFile2"), "Should not have excluded this file")
		assert.NotContains(t, files, RelativePath("dir1/nestedDir/.nestedHiddenFile1"), "Should have excluded dir1/nestedDir/**")
		assert.Contains(t, files, RelativePath("dir1/nestedDir/file3"), "Should have included this file (** should also match depth 0)")
	},

	func(t *testing.T) {
		filter := install.FileFilter{
			IgnoreHidden:        false,
			IgnoreDootCrypt:     false,
			ExploreExcludedDirs: true,
			ExcludeGlobs:        glob_collection.NewGlobCollection([]string{".hiddenDir"}),
			IncludeGlobs:        glob_collection.NewGlobCollection([]string{".hiddenDir/dïrWìthÜnicóde/1?5helloö"}),
		}
		files := install.ScanDirectory(sourceDirPath(), &filter)
		assert.Contains(t, files, RelativePath(".hiddenDir/dïrWìthÜnicóde/155helloö/fileABC"), "https://github.com/gobwas/glob/issues/54")
	},
}

func TestFileFilter_ScanDirectory(t *testing.T) {
	SetUpFiles(t, []FsNode{
		File("file1"),
		File("file2"),
		Dir("dir1", []FsNode{
			Dir("nestedDir", []FsNode{
				File("file3"),
				File(".nestedHiddenFile1"),
			}),
		}),
		File(".hiddenFile"),
		Dir(".hiddenDir", []FsNode{
			File("file4"),
			File(".nestedHiddenFile2"),
			Dir("dïrWìthÜnicóde", []FsNode{
				Dir("155helloö", []FsNode{
					File("fileABC"),
				}),
			}),
		}),
		File("secret1.doot-crypt.txt"),
		File("secret2.doot-crypt"),
		Dir("secret-dir1.doot-crypt", []FsNode{
			File("file5"),
		}),
		Dir("secret-dir2.doot-crypt.d", []FsNode{
			Dir("nested.doot-crypt", []FsNode{
				File("file6"),
				File("file7.doot-crypt"),
				File(".hiddenInSecretDir"),
			}),
		}),
	})
	for _, test := range SCAN_DIR_TESTS {
		test(t)
	}
}
