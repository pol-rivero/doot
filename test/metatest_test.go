package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pol-rivero/doot/lib/config"
)

func TestMetatest_CreateEmptyTempDirs(t *testing.T) {
	SetUp(t)
	// Check that the directories were created
	dootDir := sourceDir()
	if dootDir == "" {
		t.Fatalf("DOOT_DIR not set")
	}
	if !FileExists(dootDir) {
		t.Fatalf("DOOT_DIR does not exist")
	}

	cacheDir := cacheDir()
	if cacheDir == "" {
		t.Fatalf("DOOT_CACHE_DIR not set")
	}
	if !FileExists(cacheDir) {
		t.Fatalf("DOOT_CACHE_DIR does not exist")
	}

	targetDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Error retrieving home directory: %v", err)
	}
	if !MatchRegex(targetDir, "^/tmp/") {
		t.Fatalf("Temporary HOME was not created in /tmp (was %s)", targetDir)
	}
	if !FileExists(targetDir) {
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
	if !FileExists(dootDir, "topLevelFile") {
		t.Fatalf("topLevelFile does not exist")
	}
	fileContents, err := os.ReadFile(filepath.Join(dootDir, "topLevelFile"))
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	if string(fileContents) != "dummy text for file topLevelFile" {
		t.Fatalf("topLevelFile has unexpected contents")
	}

	if !FileExists(dootDir, "topLevelDir", "file1") {
		t.Fatalf("file1 does not exist")
	}
	if !FileExists(dootDir, "topLevelDir", "nestedDir", "file2") {
		t.Fatalf("file2 does not exist")
	}
}

func TestMetatest_CreateConfig(t *testing.T) {
	SetUp(t)
	config := config.DefaultConfig()
	config.Hosts = map[string]string{
		"my-laptop": "laptop-dots",
		"other-pc":  "other-dots",
	}
	SetUpFiles(t, []FsNode{
		ConfigFile(config),
	})
	dootDir := sourceDir()
	if !FileExists(dootDir, "config.toml") {
		t.Fatalf("config.toml does not exist")
	}
	fileContents, err := os.ReadFile(filepath.Join(dootDir, "config.toml"))
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	if !MatchRegex(string(fileContents), "^target_dir = '/tmp/TestMetatest_CreateConfig") {
		t.Fatalf("config.toml has unexpected first line: %s", string(fileContents))
	}

	if !strings.Contains(string(fileContents), "[hosts]\nmy-laptop = 'laptop-dots'\nother-pc = 'other-dots'") {
		t.Fatalf("config.toml does not contain expected hosts section: %s", string(fileContents))
	}
}
