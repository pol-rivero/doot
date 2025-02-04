package test

import (
	"testing"

	"github.com/pol-rivero/doot/lib/config"
	"github.com/pol-rivero/doot/lib/install"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/stretchr/testify/assert"
)

func TestFileMapping_SimpleMapping(t *testing.T) {
	config := config.Config{
		TargetDir:   "/target",
		ImplicitDot: false,
	}
	config.TargetDir = "/target"
	mapping := install.NewFileMapping("/src", &config, []RelativePath{
		"file1",
		".dir/.file2",
		"doot/file_that_is_internal_to_doot",
		"somedir/doot/this_is_not_the_doot_dir",
	})
	expected := []AbsolutePath{
		"/target/file1",
		"/target/.dir/.file2",
		"/target/somedir/doot/this_is_not_the_doot_dir",
	}
	assert.ElementsMatch(t, mapping.GetInstalledTargets(), expected)
}

func TestFileMapping_WithImplicitDot(t *testing.T) {
	config := config.Config{
		TargetDir:         "/target",
		ImplicitDot:       true,
		ImplicitDotIgnore: []string{"dummy_value", "not_dotted"},
	}
	config.TargetDir = "/target"
	mapping := install.NewFileMapping("/src", &config, []RelativePath{
		"file1",
		"foo/bar",
		".dir/.file2",
		"not_dotted",
		"not_dotted/file",
		"dir/not_dotted/file",
	})
	assert.ElementsMatch(t, mapping.GetInstalledTargets(), []AbsolutePath{
		"/target/.file1",
		"/target/.foo/bar",
		"/target/.dir/not_dotted/file",
		"/target/.dir/.file2",
		"/target/not_dotted",
		"/target/not_dotted/file",
	})
}
