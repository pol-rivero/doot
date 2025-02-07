package test

import (
	"os"
	"testing"

	"github.com/pol-rivero/doot/lib/add"
	"github.com/pol-rivero/doot/lib/common/config"
)

func TestAdd_BasicMapping(t *testing.T) {
	config := config.DefaultConfig()
	setUpFiles_TestAdd(t, config)
	os.Chdir(homeDir())

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
	assertHomeLink(t, "file1", sourceDir()+"/file1")
	assertHomeLink(t, "dir1/file3", sourceDir()+"/dir1/file3")
}

func TestAdd_IncorrectInputs(t *testing.T) {
	config := config.DefaultConfig()
	setUpFiles_TestAdd(t, config)
	os.Chdir(homeDir())

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
	setUpFiles_TestAdd(t, config)
	os.Chdir(homeDir() + "/dir1")

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
	assertHomeLink(t, "dir1/file3", sourceDir()+"/dir1/file3")
}

func TestAdd_WeirdInputPath(t *testing.T) {
	config := config.DefaultConfig()
	setUpFiles_TestAdd(t, config)
	os.Chdir(homeDir() + "/dir1")

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
	assertHomeLink(t, "dir1/file3", sourceDir()+"/dir1/file3")
}

func TestAdd_ExcludeInclude(t *testing.T) {
	config := config.DefaultConfig()
	config.ExcludeFiles = []string{"file1", "*.txt", "dir1", "dir3/**"}
	config.IncludeFiles = []string{"**/file6", "file2.txt"}
	setUpFiles_TestAdd(t, config)
	os.Chdir(homeDir())

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
	assertHomeLink(t, "dir3/file6", sourceDir()+"/dir3/file6")
	assertHomeLink(t, ".dir2/.foo", sourceDir()+"/.dir2/.foo")
}

func TestAdd_ImplicitDot(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = true
	config.ImplicitDotIgnore = []string{"file2.txt", "dir3"}
	setUpFiles_TestAdd(t, config)
	os.Chdir(homeDir())

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
	assertHomeLink(t, "file1", sourceDir()+"/file1")
	assertHomeLink(t, "file2.txt", sourceDir()+"/file2.txt")
	assertHomeLink(t, "dir1/file3", sourceDir()+"/dir1/file3")
}

func setUpFiles_TestAdd(t *testing.T, config config.Config) {
	SetUpFiles(t, []FsNode{
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
		}),
		Dir("dir3", []FsNode{
			File("file6"),
			File("file7"),
		}),
		Dir("emptyDir", []FsNode{}),
	}
	for _, node := range testFiles {
		createNode(homeDir(), node)
	}
}
