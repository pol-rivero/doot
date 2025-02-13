package test

import (
	"os"
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

func TestFileMapping_ConflictingNames(t *testing.T) {
	config := config.Config{
		TargetDir:   "/target",
		ImplicitDot: false,
	}
	mapping := install.NewFileMapping("/src", &config, []RelativePath{
		"file1.doot-crypt",
		"file1",
		"file2.txt",
		"file2.doot-crypt.txt",
	})
	expected := []AbsolutePath{
		"/target/file1",
		"/target/file2.txt",
	}
	assert.ElementsMatch(t, mapping.GetInstalledTargets(), expected)
}

func TestFileMapping_HostSpecificDirs(t *testing.T) {
	myHost, err := os.Hostname()
	assert.NoError(t, err)
	config := config.Config{
		TargetDir:   "/target",
		ImplicitDot: false,
		Hosts: map[string]string{
			"other_host1": "OTHER",
			myHost:        "HOST",
			"other_host2": "OTHER2",
		},
	}
	mapping := install.NewFileMapping("/src", &config, []RelativePath{
		"OTHER/file1",
		"OTHER/file-that-should-be-ignored",
		"doot/other-file-that-should-be-ignored",
		"HOST/file1",
		"file1",
		"file2",
		"OTHER2/other-ignored-file",
		"dir1/dir2/file3",
		"dir1/dir2/file4",
		"HOST/dir1/dir2/file3",
	})
	expected := []AbsolutePath{
		"/target/file1",
		"/target/file2",
		"/target/dir1/dir2/file3",
		"/target/dir1/dir2/file4",
	}
	assert.ElementsMatch(t, mapping.GetInstalledTargets(), expected)
}

func TestFileMapping_NestedHostSpecificDirs(t *testing.T) {
	myHost, err := os.Hostname()
	assert.NoError(t, err)
	config := config.Config{
		TargetDir:   "/target",
		ImplicitDot: true,
		Hosts: map[string]string{
			myHost:       "hosts/host",
			"other_host": "hosts/ignore",
		},
	}
	mapping := install.NewFileMapping("/src", &config, []RelativePath{
		"dir/file3",
		"hosts/host/file1",
		"hosts/host/.file2",
		"hosts/host/.some/nested/dir",
		"hosts/host/some/other/file",
		"hosts/ignore/file1",
		"hosts/ignore/ignore_me",
		"hosts/ignore/some/other/ignore_me",
	})
	expected := []AbsolutePath{
		"/target/.dir/file3",
		"/target/.file1",
		"/target/.file2",
		"/target/.some/nested/dir",
		"/target/.some/other/file",
	}
	assert.ElementsMatch(t, mapping.GetInstalledTargets(), expected)
}
