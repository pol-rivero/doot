package test

import (
	"os"
	"testing"

	"github.com/pol-rivero/doot/lib/commands/add"
	"github.com/pol-rivero/doot/lib/commands/restore"
	"github.com/pol-rivero/doot/lib/common/config"
	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/pol-rivero/doot/lib/utils"
	"github.com/stretchr/testify/assert"
)

func TestAdd_BasicMapping(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestAdd(t, config)
	t.Chdir(homeDir())

	assertSourceDirContents(t, "", []string{
		"doot",
	})
	add.Add([]string{
		"file1",
		homeDir() + "/dir1/file3",
	}, false, false)
	assertSourceDirContents(t, "", []string{
		"doot",
		"file1",
		"dir1",
	})
	assertSourceDirContents(t, "dir1", []string{
		"file3",
	})
	assertHomeSymlink(t, "file1", sourceDir()+"/file1")
	assertHomeSymlink(t, "dir1/file3", sourceDir()+"/dir1/file3")

	assertCache(t, []AssertCacheEntry{
		{NewAbsolutePath(homeDir() + "/file1"), sourceDir() + "/file1"},
		{NewAbsolutePath(homeDir() + "/dir1/file3"), sourceDir() + "/dir1/file3"},
	})

	assert.DirExists(t, sourceDir()+"/dir1")
	restore.Restore([]string{
		"file1",
	})
	assertHomeRegularFile(t, "file1")
	assertCache(t, []AssertCacheEntry{
		{NewAbsolutePath(homeDir() + "/dir1/file3"), sourceDir() + "/dir1/file3"},
	})
}

func TestAdd_RestoreCleansUpDirectories(t *testing.T) {
	setUpFiles_TestAdd(t, config.DefaultConfig())
	t.Chdir(homeDir())

	os.RemoveAll(sourceDir() + "/doot")
	assertSourceDirContents(t, "", []string{})

	add.Add([]string{
		".dir2/file5",
		".dir2/nested/nestedFile",
	}, false, false)
	assertSourceDirContents(t, "dir2", []string{
		"file5",
		"nested",
	})

	restore.Restore([]string{
		homeDir() + "/.dir2/file5",
		sourceDir() + "/dir2/nested/nestedFile",
	})
	assert.DirExists(t, sourceDir())
	assertSourceDirContents(t, "", []string{})
	assertCache(t, []AssertCacheEntry{})
}

func TestAdd_IncorrectInputs(t *testing.T) {
	config := config.DefaultConfig()
	setUpFiles_TestAdd(t, config)
	t.Chdir(homeDir())

	add.Add([]string{
		"some-nonexistent-file", // Doesn't exist
		"dir1",                  // Is a directory
		"/etc/passwd",           // Outside target (home) directory
	}, false, false)
	assertSourceDirContents(t, "", []string{
		"doot",
	})
}

func TestAdd_FromAnotherPWD(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestAdd(t, config)
	t.Chdir(homeDir() + "/dir1")

	add.Add([]string{
		"file3",
		"../file2.txt",
	}, false, false)
	assertSourceDirContents(t, "", []string{
		"doot",
		"file2.txt",
		"dir1",
	})
	assertSourceDirContents(t, "dir1", []string{
		"file3",
	})
	assertHomeSymlink(t, "dir1/file3", sourceDir()+"/dir1/file3")
}

func TestAdd_WeirdInputPath(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestAdd(t, config)
	t.Chdir(homeDir() + "/dir1")

	add.Add([]string{
		"../dir3////./../dir1//file3",
	}, false, false)
	assertSourceDirContents(t, "", []string{
		"doot",
		"dir1",
	})
	assertSourceDirContents(t, "dir1", []string{
		"file3",
	})
	assertHomeSymlink(t, "dir1/file3", sourceDir()+"/dir1/file3")
}

func TestAdd_ExcludeInclude(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	config.ExcludeFiles = []string{"file1", "*.txt", "dir1", "dir3/**"}
	config.IncludeFiles = []string{"**/file6", "file2.txt"}
	setUpFiles_TestAdd(t, config)
	t.Chdir(homeDir())

	add.Add([]string{
		"file1",                // Excluded
		"file2.txt",            // Excluded but then included
		"dir1/nestedDir/file4", // Parent dir excluded
		"dir3/file6",           // Excluded but then included
		"dir3/file7",           // Excluded
		".dir2/.foo",           // Not excluded
	}, false, false)
	assertSourceDirContents(t, "", []string{
		"doot",
		"file2.txt",
		"dir3",
		".dir2",
	})
	assertHomeSymlink(t, "dir3/file6", sourceDir()+"/dir3/file6")
	assertHomeSymlink(t, ".dir2/.foo", sourceDir()+"/.dir2/.foo")
}

func TestAdd_ImplicitDot(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = true
	config.ImplicitDotIgnore = []string{"file2.txt", "dir3"}
	config.ExcludeFiles = []string{}
	setUpFiles_TestAdd(t, config)
	t.Chdir(homeDir())

	add.Add([]string{
		"file1", // Impossible filename with implicit dot
		"file2.txt",
		"dir1/file3", // Impossible filename with implicit dot
		".dir2/file5",
		".dir2/.foo",
		"dir3/file6",
	}, false, false)
	assertSourceDirContents(t, "", []string{
		"doot",
		"file2.txt",
		"dir2",
		"dir3",
	})
	assertSourceDirContents(t, "dir2", []string{
		"file5",
		".foo",
	})
	assertHomeSymlink(t, "file2.txt", sourceDir()+"/file2.txt")
	assertHomeSymlink(t, ".dir2/file5", sourceDir()+"/dir2/file5")
	assertHomeSymlink(t, "dir3/file6", sourceDir()+"/dir3/file6")
}

func TestAdd_WithCryptExtensionUninitialized(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestAdd(t, config)
	t.Chdir(homeDir())

	add.Add([]string{
		"file1",
		"file2.txt",
		"dir1/nestedDir/file4",
		"dir.with.dots/file.with.some.dots",
		"dir.with.dots/file-without-dots",
	}, true, false)

	assertSourceDirContents(t, "", []string{
		"doot",
	})
}

func TestAdd_WithCryptExtension(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	config.ExcludeFiles = []string{}
	setUpFiles_TestAdd(t, config)
	initializeGitCrypt()
	t.Chdir(homeDir())

	add.Add([]string{
		"file1",
		"file2.txt",
		"dir1/nestedDir/file4",
		"dir.with.dots/file.with.some.dots",
		"dir.with.dots/file-without-dots",
		".dir2/.foo",
		".dir2/.foo.txt",
	}, true, false)
	assertSourceDirContents(t, "", []string{
		".git",
		"doot",
		"file1.doot-crypt",
		"file2.doot-crypt.txt",
		"dir1",
		"dir.with.dots",
		".dir2",
	})
	assertSourceDirContents(t, "dir1", []string{
		"nestedDir",
	})
	assertSourceDirContents(t, "dir1/nestedDir", []string{
		"file4.doot-crypt",
	})
	assertSourceDirContents(t, "dir.with.dots", []string{
		"file.with.some.doot-crypt.dots",
		"file-without-dots.doot-crypt",
	})
	assertSourceDirContents(t, ".dir2", []string{
		".foo.doot-crypt",
		".foo.doot-crypt.txt",
	})
	assertHomeSymlink(t, "file1", sourceDir()+"/file1.doot-crypt")
	assertHomeSymlink(t, "file2.txt", sourceDir()+"/file2.doot-crypt.txt")
	assertHomeSymlink(t, "dir1/nestedDir/file4", sourceDir()+"/dir1/nestedDir/file4.doot-crypt")
	assertHomeSymlink(t, "dir.with.dots/file.with.some.dots", sourceDir()+"/dir.with.dots/file.with.some.doot-crypt.dots")
	assertHomeSymlink(t, "dir.with.dots/file-without-dots", sourceDir()+"/dir.with.dots/file-without-dots.doot-crypt")
	assertHomeSymlink(t, ".dir2/.foo", sourceDir()+"/.dir2/.foo.doot-crypt")
	assertHomeSymlink(t, ".dir2/.foo.txt", sourceDir()+"/.dir2/.foo.doot-crypt.txt")
}

func TestAdd_ExcludeIncludeWithCrypt(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	config.ExcludeFiles = []string{"file1.doot-crypt", "*.txt", "dir1", "dir3/**"}
	config.IncludeFiles = []string{"**/file6.doot-crypt", "file2.doot-crypt.txt"}
	setUpFiles_TestAdd(t, config)
	initializeGitCrypt()
	t.Chdir(homeDir())

	add.Add([]string{
		"file1",                // Excluded
		"file2.txt",            // Excluded but then included
		"dir1/nestedDir/file4", // Parent dir excluded
		"dir3/file6",           // Excluded but then included
		"dir3/file7",           // Excluded
		".dir2/.foo",           // Not excluded
	}, true, false)
	assertSourceDirContents(t, "", []string{
		".git",
		"doot",
		"file2.doot-crypt.txt",
		"dir3",
		".dir2",
	})
	assertHomeSymlink(t, "dir3/file6", sourceDir()+"/dir3/file6.doot-crypt")
	assertHomeSymlink(t, ".dir2/.foo", sourceDir()+"/.dir2/.foo.doot-crypt")
}

func TestAdd_WithCryptDirectory(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = true
	config.ImplicitDotIgnore = []string{
		"cryptTest",
		"cryptTest.doot-crypt",
	}
	config.ExcludeFiles = []string{}
	setUpFiles_TestAdd(t, config)
	initializeGitCrypt()
	t.Chdir(homeDir())

	createNode(sourceDir(), Dir("cryptTest.doot-crypt", []FsNode{}))
	createNode(sourceDir(), Dir("cryptTest", []FsNode{
		Dir("foo.doot-crypt", []FsNode{}),
	}))

	// Prefers cryptTest over cryptTest.doot-crypt. Uses foo.doot-crypt because foo doesn't exist
	// Also, doesn't add .doot-crypt to the name because the path already contains .doot-crypt
	add.Add([]string{
		"cryptTest/foo/secret1.txt",
	}, true, false)
	assertHomeSymlink(t, "cryptTest/foo/secret1.txt", sourceDir()+"/cryptTest/foo.doot-crypt/secret1.txt")

	os.RemoveAll(sourceDir() + "/cryptTest")
	// Now there's no choice but to use the cryptTest.doot-crypt directory and create 'foo'
	add.Add([]string{
		"cryptTest/foo/secret2.txt",
	}, true, false)
	assertHomeSymlink(t, "cryptTest/foo/secret2.txt", sourceDir()+"/cryptTest.doot-crypt/foo/secret2.txt")

	utils.USER_INPUT_MOCK_RESPONSE = "y"
	add.Add([]string{
		"cryptTest/foo/not-really-a-secret1.txt",
	}, false, false)
	assertHomeSymlink(t, "cryptTest/foo/not-really-a-secret1.txt", sourceDir()+"/cryptTest.doot-crypt/foo/not-really-a-secret1.txt")

	utils.USER_INPUT_MOCK_RESPONSE = "n"
	add.Add([]string{
		"cryptTest/foo/not-really-a-secret2.txt",
	}, false, false)
	assertHomeSymlink(t, "cryptTest/foo/not-really-a-secret2.txt", sourceDir()+"/cryptTest/foo/not-really-a-secret2.txt")
}

func TestAdd_HostSpecificNotFound(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	config.Hosts = map[string]string{
		"other-host": "foo",
	}
	setUpFiles_TestAdd(t, config)
	t.Chdir(homeDir())

	log.PanicInsteadOfExit = true
	assert.Panics(t, func() {
		add.Add([]string{"file1"}, false, true)
	})

	assertSourceDirContents(t, "", []string{
		"doot",
	})
}

func TestAdd_HostSpecificDir(t *testing.T) {
	host, err := os.Hostname()
	assert.NoError(t, err)

	config := config.DefaultConfig()
	config.ImplicitDot = false
	config.Hosts = map[string]string{
		"other-host": "foo",
		host:         "host/dir",
	}
	setUpFiles_TestAdd(t, config)
	t.Chdir(homeDir())

	add.Add([]string{
		"file2.txt",
		"dir1/nestedDir/file4",
	}, false, true)
	assertSourceDirContents(t, "", []string{
		"doot",
		"host",
	})
	assertSourceDirContents(t, "host/dir", []string{
		"file2.txt",
		"dir1",
	})
	assertSourceDirContents(t, "host/dir/dir1/nestedDir", []string{
		"file4",
	})
	assertHomeSymlink(t, "file2.txt", sourceDir()+"/host/dir/file2.txt")
	assertHomeSymlink(t, "dir1/nestedDir/file4", sourceDir()+"/host/dir/dir1/nestedDir/file4")
}

func TestAdd_IsIdempotent(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestAdd(t, config)
	t.Chdir(homeDir())

	add.Add([]string{
		"file1",
		"file2.txt",
		"dir1/file3",
	}, false, false)

	os.Remove(homeDir() + "/file1")
	createNode(homeDir(), FsFile{Name: "file1", Content: "new content"})
	os.Remove(homeDir() + "/file2.txt")
	createSymlink(homeDir(), "file2.txt", "./file1")

	utils.USER_INPUT_MOCK_RESPONSE = "n"
	add.Add([]string{
		"file1",      // Dotfile already exists and this is a regular file, fails with error
		"file2.txt",  // Dotfile already exists and this points to wrong location, fails with error
		"dir1/file3", // Already points to the same source, skipped
	}, false, false)

	assertSourceDirContents(t, "", []string{
		"doot",
		"file1",
		"file2.txt",
		"dir1",
	})
	assertHomeRegularFile(t, "file1")
	assertHomeSymlink(t, "file2.txt", "./file1")
	assertHomeSymlink(t, "dir1/file3", sourceDir()+"/dir1/file3")
}

func setUpFiles_TestAdd(t *testing.T, config config.Config) {
	SetUpFiles(t, true, []FsNode{
		Dir("doot", []FsNode{
			ConfigFile(config),
		}),
	})
	testFiles := []FsNode{
		File("file1"),
		File("file2.txt"),
		Dir("dir1", []FsNode{
			File(".file2"),
			File("file3"),
			Dir("nestedDir", []FsNode{
				File("file4"),
			}),
		}),
		Dir(".dir2", []FsNode{
			File("file5"),
			File(".foo"),
			File(".foo.txt"),
			Dir("nested", []FsNode{
				File("nestedFile"),
			}),
		}),
		Dir("dir3", []FsNode{
			File("file6"),
			File("file7"),
		}),
		Dir("dir.with.dots", []FsNode{
			File("file.with.some.dots"),
			File("file-without-dots"),
		}),
		Dir("emptyDir", []FsNode{}),
		Dir("cryptTest", []FsNode{
			Dir("foo", []FsNode{
				File("secret1.txt"),
				File("secret2.txt"),
				File("not-really-a-secret1.txt"),
				File("not-really-a-secret2.txt"),
			}),
		}),
	}
	for _, node := range testFiles {
		createNode(homeDir(), node)
	}
}
