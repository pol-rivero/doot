package test

import (
	"testing"

	"github.com/pol-rivero/doot/lib/common/config"
	"github.com/pol-rivero/doot/lib/common/log"
	"github.com/stretchr/testify/assert"
)

func loadConfigWithHosts(t *testing.T, hosts map[string]string) config.Config {
	cfg := config.DefaultConfig()
	cfg.TargetDir = "$HOME"
	cfg.Hosts = hosts
	SetUpFiles(t, false, []FsNode{
		Dir("doot", []FsNode{
			ConfigFile(cfg),
		}),
	})
	return config.FromDotfilesDir(sourceDirPath())
}

func TestConfig_InvalidHostDirs(t *testing.T) {
	invalidDirs := []string{
		"",
		".",
		"./",
		"..",
		"../escaped",
		"foo/../../escaped",
		"/absolute/path",
		"doot",
		"doot/foo",
	}
	log.PanicInsteadOfExit = true
	for _, dir := range invalidDirs {
		assert.Panics(t, func() {
			loadConfigWithHosts(t, map[string]string{"some-host": dir})
		}, "Host dir '%s' should be rejected", dir)
	}
}

func TestConfig_HostDirsAreCleaned(t *testing.T) {
	cfg := loadConfigWithHosts(t, map[string]string{
		"host1": "laptop/",
		"host2": "./foo//bar/",
		"host3": "foo/../baz",
		"host4": "doot-files",
	})
	assert.Equal(t, map[string]string{
		"host1": "laptop",
		"host2": "foo/bar",
		"host3": "baz",
		"host4": "doot-files",
	}, cfg.Hosts)
}

func TestConfig_ParseErrorIsFatal(t *testing.T) {
	invalidConfigs := []string{
		"target_dir = \"$HOME\"\nimplicit_dot = false\n[hosts",
		"exclude_files = [\"README.md\"",
		"implicit_dot = notABool",
		"use_hardlinks = \"yes\"",
	}
	log.PanicInsteadOfExit = true
	for _, content := range invalidConfigs {
		SetUpFiles(t, false, []FsNode{
			Dir("doot", []FsNode{
				FsFile{Name: "config.toml", Content: content},
			}),
		})
		assert.Panics(t, func() {
			config.FromDotfilesDir(sourceDirPath())
		}, "Config '%s' should be rejected", content)
	}
}
