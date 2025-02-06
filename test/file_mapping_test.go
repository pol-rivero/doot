package test

import (
	"testing"

	"github.com/pol-rivero/doot/lib/common/config"
	"github.com/pol-rivero/doot/lib/install"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/stretchr/testify/assert"
)

func TestFileMapping_SimpleMapping(t *testing.T) {
	config := config.Config{
		TargetDir:   "/target",
		ImplicitDot: false,
	}
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

func TestFileMapping_WithDootCrypt(t *testing.T) {
	config := config.Config{
		TargetDir:   "/target",
		ImplicitDot: false,
	}
	mapping := install.NewFileMapping("/src", &config, []RelativePath{
		"file1.doot-crypt",
		"file2.doot-crypt.txt",
		"dirA.doot-crypt/dirB.doot-crypt.d/file3.doot-crypt",
	})
	expected := []AbsolutePath{
		"/target/file1",
		"/target/file2.txt",
		"/target/dirA/dirB.d/file3",
	}
	assert.ElementsMatch(t, mapping.GetInstalledTargets(), expected)
}
