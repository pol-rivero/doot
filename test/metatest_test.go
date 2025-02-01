package test

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMetatest_CreateEmptyTempDirs(t *testing.T) {
	SetUp(t)
	// Check that the directories were created
	dootDir := sourceDir()
	if dootDir == "" {
		t.Fatalf("DOOT_DIR not set")
	}
	if _, err := os.Stat(dootDir); err != nil {
		t.Fatalf("DOOT_DIR does not exist")
	}

	cacheDir := cacheDir()
	if cacheDir == "" {
		t.Fatalf("DOOT_CACHE_DIR not set")
	}
	if _, err := os.Stat(cacheDir); err != nil {
		t.Fatalf("DOOT_CACHE_DIR does not exist")
	}

	targetDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Error retrieving home directory: %v", err)
	}
	if !MatchRegex(targetDir, "^/tmp/") {
		t.Fatalf("Temporary HOME was not created in /tmp (was %s)", targetDir)
	}
	if _, err := os.Stat(targetDir); err != nil {
		t.Fatalf("HOME does not exist")
	}

	homeDir := homeDir()
	if homeDir != targetDir {
		t.Fatalf("homeDir() returned %s, expected %s", homeDir, targetDir)
	}
}

func TestMetatest_CreateTempDirs(t *testing.T) {
	SetUpFiles(t, []FsNode{
		File("topLevelFile"),
		Dir("topLevelDir", []FsNode{
			File("file1"),
			Dir("nestedDir", []FsNode{
				File("file2"),
			}),
		}),
	})
	dootDir := sourceDir()
	if _, err := os.Stat(filepath.Join(dootDir, "topLevelFile")); err != nil {
		t.Fatalf("topLevelFile does not exist")
	}
	fileContents, err := os.ReadFile(filepath.Join(dootDir, "topLevelFile"))
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	if !MatchRegex(string(fileContents), `^dummy text for file .*/topLevelFile$`) {
		t.Fatalf("topLevelFile has unexpected contents")
	}

	if _, err := os.Stat(filepath.Join(dootDir, "topLevelDir", "file1")); err != nil {
		t.Fatalf("file1 does not exist")
	}
	if _, err := os.Stat(filepath.Join(dootDir, "topLevelDir", "nestedDir", "file2")); err != nil {
		t.Fatalf("file2 does not exist")
	}
}
