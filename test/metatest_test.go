package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pol-rivero/doot/lib/common/config"
	"github.com/stretchr/testify/assert"
)

func TestMetatest_CreateEmptyTempDirs(t *testing.T) {
	SetUp(t)
	dootDir := sourceDir()
	assert.NotEmpty(t, dootDir, "DOOT_DIR not set")
	assert.DirExists(t, dootDir, "DOOT_DIR does not exist")

	cacheDir := cacheDir()
	assert.NotEmpty(t, cacheDir, "DOOT_CACHE_DIR not set")
	assert.DirExists(t, cacheDir, "DOOT_CACHE_DIR does not exist")

	targetDir, err := os.UserHomeDir()
	assert.NoError(t, err, "Error retrieving home directory")
	assert.Regexp(t, "^/tmp/", targetDir, "Temporary HOME was not created in /tmp")
	assert.DirExists(t, targetDir, "HOME does not exist")

	homeDir := homeDir()
	assert.Equal(t, targetDir, homeDir, "homeDir() returned unexpected value")
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
	assert.FileExists(t, dootDir+"/topLevelFile")

	fileContents, err := os.ReadFile(dootDir + "/topLevelFile")
	assert.NoError(t, err, "Error reading file")
	assert.Equal(t, "dummy text for file topLevelFile", string(fileContents), "topLevelFile has unexpected contents")

	assert.FileExists(t, dootDir+"/topLevelDir/file1")
	assert.FileExists(t, dootDir+"/topLevelDir/nestedDir/file2")
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
	assert.FileExists(t, filepath.Join(dootDir, "config.toml"))

	fileContents, err := os.ReadFile(filepath.Join(dootDir, "config.toml"))
	assert.NoError(t, err, "Error reading file")
	assert.Regexp(t, "^target_dir = '\\$HOME", string(fileContents), "config.toml has unexpected first line")
	assert.Contains(t, string(fileContents), "[hosts]\nmy-laptop = 'laptop-dots'\nother-pc = 'other-dots'", "config.toml does not contain expected hosts section")
}

func TestMetatest_CreateConfigBeforeSetUp(t *testing.T) {
	config := config.DefaultConfig()
	SetUpFiles(t, []FsNode{
		ConfigFile(config),
	})
	dootDir := sourceDir()
	assert.FileExists(t, filepath.Join(dootDir, "config.toml"))

	fileContents, err := os.ReadFile(filepath.Join(dootDir, "config.toml"))
	assert.NoError(t, err, "Error reading file")
	assert.Regexp(t, "^target_dir = '\\$HOME", string(fileContents), "$HOME was not replaced before writing config.toml")
}
