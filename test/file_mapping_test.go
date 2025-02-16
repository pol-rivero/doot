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
	assertSymlinkCollection(t, mapping.GetInstalledTargets(), map[AbsolutePath]AbsolutePath{
		"/target/file1":       "/src/file1",
		"/target/.dir/.file2": "/src/.dir/.file2",
		"/target/somedir/doot/this_is_not_the_doot_dir": "/src/somedir/doot/this_is_not_the_doot_dir",
	})
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
	assertSymlinkCollection(t, mapping.GetInstalledTargets(), map[AbsolutePath]AbsolutePath{
		"/target/.file1":               "/src/file1",
		"/target/.foo/bar":             "/src/foo/bar",
		"/target/.dir/not_dotted/file": "/src/dir/not_dotted/file",
		"/target/.dir/.file2":          "/src/.dir/.file2",
		"/target/not_dotted":           "/src/not_dotted",
		"/target/not_dotted/file":      "/src/not_dotted/file",
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
	assertSymlinkCollection(t, mapping.GetInstalledTargets(), map[AbsolutePath]AbsolutePath{
		"/target/file1":             "/src/file1.doot-crypt",
		"/target/file2.txt":         "/src/file2.doot-crypt.txt",
		"/target/dirA/dirB.d/file3": "/src/dirA.doot-crypt/dirB.doot-crypt.d/file3.doot-crypt",
	})
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
	assertSymlinkCollection(t, mapping.GetInstalledTargets(), map[AbsolutePath]AbsolutePath{
		"/target/file1":     "/src/file1.doot-crypt", // Prefer the first file found
		"/target/file2.txt": "/src/file2.txt",
	})
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
	assertSymlinkCollection(t, mapping.GetInstalledTargets(), map[AbsolutePath]AbsolutePath{
		"/target/file1":           "/src/HOST/file1",
		"/target/file2":           "/src/file2",
		"/target/dir1/dir2/file3": "/src/HOST/dir1/dir2/file3",
		"/target/dir1/dir2/file4": "/src/dir1/dir2/file4",
	})
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
	assertSymlinkCollection(t, mapping.GetInstalledTargets(), map[AbsolutePath]AbsolutePath{
		"/target/.dir/file3":       "/src/dir/file3",
		"/target/.file1":           "/src/hosts/host/file1",
		"/target/.file2":           "/src/hosts/host/.file2",
		"/target/.some/nested/dir": "/src/hosts/host/.some/nested/dir",
		"/target/.some/other/file": "/src/hosts/host/some/other/file",
	})
}
