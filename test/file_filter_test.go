package test

import (
	"slices"
	"testing"

	"github.com/pol-rivero/doot/lib/config"
	"github.com/pol-rivero/doot/lib/helpers"
)

func TestFileFilter_CreateFilter1(t *testing.T) {
	config := &config.Config{
		ExcludeFiles: []string{"**/.*", "file1"},
		IncludeFiles: []string{"file2"},
	}
	filter := helpers.CreateFilter(config, false)
	if filter.IgnoreHidden != true {
		t.Fatalf("Expected IgnoreHidden to be true")
	}
	if filter.IgnoreDootCrypt != false {
		t.Fatalf("Expected IgnoreDootCrypt to be false")
	}
	if filter.ExcludeGlobs.Len() != 1 {
		t.Fatalf("Expected 1 exclude glob, got %d", filter.ExcludeGlobs.Len())
	}
	if filter.IncludeGlobs.Len() != 1 {
		t.Fatalf("Expected 1 include glob, got %d", filter.IncludeGlobs.Len())
	}
}

func TestFileFilter_CreateFilter2(t *testing.T) {
	invalid_glob := "*["
	config := &config.Config{
		ExcludeFiles: []string{"file1", "*.txt"},
		IncludeFiles: []string{"file2", invalid_glob},
	}
	filter := helpers.CreateFilter(config, true)
	if filter.IgnoreHidden != false {
		t.Fatalf("Expected IgnoreHidden to be false")
	}
	if filter.IgnoreDootCrypt != true {
		t.Fatalf("Expected IgnoreDootCrypt to be true")
	}
	if filter.ExcludeGlobs.Len() != 2 {
		t.Fatalf("Expected 2 exclude globs, got %d", filter.ExcludeGlobs.Len())
	}
	if filter.IncludeGlobs.Len() != 1 {
		t.Fatalf("Expected 1 include glob (the other one is invalid), got %d", filter.IncludeGlobs.Len())
	}
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
	scanAll(t)
	ignoreHidden(t)
	ignoreCrypt(t)
	ignoreHiddenAndCrypt(t)
	excludeAndInclude1(t)
	excludeAndInclude2(t)
	excludeAndInclude3(t)
	weirdSuperAsterisk(t)
}

func scanAll(t *testing.T) {
	filter := helpers.FileFilter{
		IgnoreHidden:    false,
		IgnoreDootCrypt: false,
		ExcludeGlobs:    helpers.NewGlobCollection([]string{}),
		IncludeGlobs:    helpers.NewGlobCollection([]string{}),
	}
	files := helpers.ScanDirectory(sourceDir(), filter)
	if !slices.Equal(files, []string{
		".hiddenDir/.nestedHiddenFile2",
		".hiddenDir/file4",
		".hiddenFile",
		"dir1/nestedDir/.nestedHiddenFile1",
		"dir1/nestedDir/file3",
		"file1",
		"file2",
		"secret-dir1.doot-crypt/file5",
		"secret-dir2.doot-crypt.d/nested.doot-crypt/.hiddenInSecretDir",
		"secret-dir2.doot-crypt.d/nested.doot-crypt/file6",
		"secret-dir2.doot-crypt.d/nested.doot-crypt/file7.doot-crypt",
		"secret1.doot-crypt.txt",
		"secret2.doot-crypt",
	}) {
		t.Fatalf("Unexpected files: %v", files)
	}
}

func ignoreHidden(t *testing.T) {
	filter := helpers.FileFilter{
		IgnoreHidden:    true,
		IgnoreDootCrypt: false,
		ExcludeGlobs:    helpers.NewGlobCollection([]string{}),
		IncludeGlobs:    helpers.NewGlobCollection([]string{}),
	}
	files := helpers.ScanDirectory(sourceDir(), filter)
	if !slices.Equal(files, []string{
		"dir1/nestedDir/file3",
		"file1",
		"file2",
		"secret-dir1.doot-crypt/file5",
		"secret-dir2.doot-crypt.d/nested.doot-crypt/file6",
		"secret-dir2.doot-crypt.d/nested.doot-crypt/file7.doot-crypt",
		"secret1.doot-crypt.txt",
		"secret2.doot-crypt",
	}) {
		t.Fatalf("Unexpected files: %v", files)
	}
}

func ignoreCrypt(t *testing.T) {
	filter := helpers.FileFilter{
		IgnoreHidden:    false,
		IgnoreDootCrypt: true,
		ExcludeGlobs:    helpers.NewGlobCollection([]string{}),
		IncludeGlobs:    helpers.NewGlobCollection([]string{}),
	}
	files := helpers.ScanDirectory(sourceDir(), filter)
	if !slices.Equal(files, []string{
		".hiddenDir/.nestedHiddenFile2",
		".hiddenDir/file4",
		".hiddenFile",
		"dir1/nestedDir/.nestedHiddenFile1",
		"dir1/nestedDir/file3",
		"file1",
		"file2",
	}) {
		t.Fatalf("Unexpected files: %v", files)
	}
}

func ignoreHiddenAndCrypt(t *testing.T) {
	filter := helpers.FileFilter{
		IgnoreHidden:    true,
		IgnoreDootCrypt: true,
		ExcludeGlobs:    helpers.NewGlobCollection([]string{}),
		IncludeGlobs:    helpers.NewGlobCollection([]string{}),
	}
	files := helpers.ScanDirectory(sourceDir(), filter)
	if !slices.Equal(files, []string{
		"dir1/nestedDir/file3",
		"file1",
		"file2",
	}) {
		t.Fatalf("Unexpected files: %v", files)
	}
}

func excludeAndInclude1(t *testing.T) {
	filter := helpers.FileFilter{
		IgnoreHidden:    false,
		IgnoreDootCrypt: false,
		ExcludeGlobs:    helpers.NewGlobCollection([]string{"secret*"}),
		IncludeGlobs:    helpers.NewGlobCollection([]string{"*.txt", "**/file6"}),
	}
	files := helpers.ScanDirectory(sourceDir(), filter)
	if slices.Contains(files, "secret2.doot-crypt") {
		t.Fatalf("Should have excluded secret2.doot-crypt")
	}
	if !slices.Contains(files, "secret1.doot-crypt.txt") {
		t.Fatalf("Should have included secret1.doot-crypt.txt because it ends with .txt")
	}
	if slices.Contains(files, "secret-dir2.doot-crypt.d/nested.doot-crypt/file6") {
		t.Fatalf("Even though file6 is included, it shouldn't be returned because excluded directories are not explored")
	}
}

func excludeAndInclude2(t *testing.T) {
	filter := helpers.FileFilter{
		IgnoreHidden:    false,
		IgnoreDootCrypt: false,
		ExcludeGlobs:    helpers.NewGlobCollection([]string{"secret*/**"}),
		IncludeGlobs:    helpers.NewGlobCollection([]string{"secret*/nested.doot-crypt", "**/file6"}),
	}
	files := helpers.ScanDirectory(sourceDir(), filter)
	if slices.Contains(files, "secret-dir1.doot-crypt/file5") {
		t.Fatalf("Should have excluded secret-dir1.doot-crypt/file5")
	}
	if !slices.Contains(files, "secret-dir2.doot-crypt.d/nested.doot-crypt/file6") {
		t.Fatalf("Now file6 should be returned, because the directory contents are excluded, not the directory itself")
	}
	if slices.Contains(files, "secret-dir2.doot-crypt.d/nested.doot-crypt/file7.doot-crypt") {
		t.Fatalf("Should have excluded secret-dir2.doot-crypt.d/nested.doot-crypt/file7.doot-crypt")
	}
}

func excludeAndInclude3(t *testing.T) {
	filter := helpers.FileFilter{
		IgnoreHidden:    false,
		IgnoreDootCrypt: false,
		ExcludeGlobs:    helpers.NewGlobCollection([]string{"secret*/**"}),
		IncludeGlobs:    helpers.NewGlobCollection([]string{"secret*/nested.doot-crypt**"}),
	}
	files := helpers.ScanDirectory(sourceDir(), filter)
	if slices.Contains(files, "secret-dir1.doot-crypt/file5") {
		t.Fatalf("Should have excluded secret-dir1.doot-crypt/file5")
	}
	if !slices.Contains(files, "secret-dir2.doot-crypt.d/nested.doot-crypt/file6") {
		t.Fatalf("All children of nested.doot-crypt should be included")
	}
	if !slices.Contains(files, "secret-dir2.doot-crypt.d/nested.doot-crypt/file7.doot-crypt") {
		t.Fatalf("All children of nested.doot-crypt should be included")
	}
}

func weirdSuperAsterisk(t *testing.T) {
	filter := helpers.FileFilter{
		IgnoreHidden:    false,
		IgnoreDootCrypt: false,
		ExcludeGlobs:    helpers.NewGlobCollection([]string{"**/file2", ".hiddenDir/**/file4", "dir1/nestedDir/**"}),
		IncludeGlobs:    helpers.NewGlobCollection([]string{"dir1/nestedDir/**/file3"}),
	}
	files := helpers.ScanDirectory(sourceDir(), filter)
	if !slices.Contains(files, "file1") {
		t.Fatalf("Should not have excluded file1")
	}
	if slices.Contains(files, "file2") {
		// t.Fatalf("Should have excluded file2")
		t.Log("This is a bug in gobwas/glob, see https://github.com/gobwas/glob/issues/58")
	}
	if slices.Contains(files, ".hiddenDir/file4") {
		t.Fatalf("Should have excluded .hiddenDir/file4")
	}
	if !slices.Contains(files, ".hiddenDir/.nestedHiddenFile2") {
		t.Fatalf("Should not have excluded .hiddenDir/.nestedHiddenFile2")
	}
	if slices.Contains(files, "dir1/nestedDir/.nestedHiddenFile1") {
		t.Fatalf("Should have excluded dir1/nestedDir/.nestedHiddenFile1")
	}
	if !slices.Contains(files, "dir1/nestedDir/file3") {
		t.Fatalf("Should have included dir1/nestedDir/file3")
	}
}
