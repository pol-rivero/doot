package test

import (
	"os"
	"testing"

	"github.com/pol-rivero/doot/lib/common/cache"
	"github.com/pol-rivero/doot/lib/common/config"
	"github.com/pol-rivero/doot/lib/install"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/pol-rivero/doot/lib/utils"
	"github.com/stretchr/testify/assert"
)

func TestInstall_DefaultConfig(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestInstall(t, config)

	install.Install()
	assertHomeDirContents(t, "", []string{
		"file1",
		"file2.txt",
		"dir1",
		"dir3",
	})
	assertHomeDirContents(t, "dir1", []string{
		"file3",
		"nestedDir",
	})
	assertHomeDirContents(t, "dir1/nestedDir", []string{
		"file4",
	})
	assertHomeDirContents(t, "dir3", []string{
		"file6",
	})
	assertHomeLink(t, "file1", sourceDir()+"/file1")
	assertHomeLink(t, "dir1/nestedDir/file4", sourceDir()+"/dir1/nestedDir/file4")

	install.Clean()
	assertHomeDirContents(t, "", []string{})
}

func TestInstall_HiddenFiles(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	config.ExcludeFiles = []string{"file2.txt"}
	setUpFiles_TestInstall(t, config)

	install.Install()
	assertHomeDirContents(t, "", []string{
		"file1",
		"dir1",
		".dir2",
		"dir3",
	})
	assertHomeDirContents(t, "dir1", []string{
		".file2",
		"file3",
		"nestedDir",
	})
	assertHomeLink(t, ".dir2/file5", sourceDir()+"/.dir2/file5")

	install.Clean()
	assertHomeDirContents(t, "", []string{})
}

func TestInstall_ImplicitDot(t *testing.T) {
	config := config.DefaultConfig()
	config.ExcludeFiles = []string{}
	config.ImplicitDot = true
	config.ImplicitDotIgnore = []string{"file2.txt", "dir3"}
	setUpFiles_TestInstall(t, config)

	install.Install()
	assertHomeDirContents(t, "", []string{
		".file1",
		"file2.txt",
		".dir1",
		".dir2",
		"dir3",
	})
	assertHomeDirContents(t, ".dir1", []string{
		".file2",
		"file3",
		"nestedDir",
	})
	assertHomeDirContents(t, "dir3", []string{
		"file6",
	})
	assertHomeLink(t, ".file1", sourceDir()+"/file1")
	assertHomeLink(t, ".dir1/file3", sourceDir()+"/dir1/file3")

	install.Clean()
	assertHomeDirContents(t, "", []string{})
}

func TestInstall_MixedWithRegularFiles(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestInstall(t, config)
	createFile(homeDir(), File("existingFile"))
	createDir(homeDir(), Dir("dir1", []FsNode{
		File("existingFileInDir1"),
	}))
	assertHomeDirContents(t, "", []string{
		"existingFile",
		"dir1",
	})

	install.Install()
	assertHomeDirContents(t, "", []string{
		"existingFile",
		"file1",
		"file2.txt",
		"dir1",
		"dir3",
	})
	assertHomeDirContents(t, "dir1", []string{
		"existingFileInDir1",
		"file3",
		"nestedDir",
	})

	install.Clean()
	assertHomeDirContents(t, "", []string{
		"existingFile",
		"dir1",
	})
	assertHomeDirContents(t, "dir1", []string{
		"existingFileInDir1",
	})
}

func TestInstall_UpdatesCache(t *testing.T) {
	config := config.DefaultConfig()
	config.ExcludeFiles = []string{}
	config.ImplicitDot = true
	config.ImplicitDotIgnore = []string{"file2.txt", "dir3"}

	setUpFiles_TestInstall(t, config)
	install.Install()

	homePath := NewAbsolutePath(homeDir())
	assertCache(t, []AssertCacheEntry{
		{Path: homePath.Join(".file1"), Content: sourceDir() + "/file1"},
		{Path: homePath.Join("file2.txt"), Content: sourceDir() + "/file2.txt"},
		{Path: homePath.Join(".dir1/.file2"), Content: sourceDir() + "/dir1/.file2"},
		{Path: homePath.Join(".dir1/file3"), Content: sourceDir() + "/dir1/file3"},
		{Path: homePath.Join(".dir1/nestedDir/file4"), Content: sourceDir() + "/dir1/nestedDir/file4"},
		{Path: homePath.Join(".dir2/file5"), Content: sourceDir() + "/.dir2/file5"},
		{Path: homePath.Join("dir3/file6"), Content: sourceDir() + "/dir3/file6"},
	})

	install.Clean()
	assertCache(t, []AssertCacheEntry{})
}

func TestInstall_IncrementalInstall(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestInstall(t, config)
	createSymlink(homeDir(), "someFileInstalledInAPreviousRun", sourceDir()+"/file1")

	// Lie to the cache and see that only the other files were added.
	// Also, someFileInstalledInAPreviousRun is no longer in dotfiles dir, so it should be removed
	dootCache := cache.Load()
	cacheEntry := dootCache.GetEntry(sourceDir() + ":" + homeDir())
	cacheEntry.Links = []*cache.InstalledFile{
		{
			Path:    homeDir() + "/file1",
			Content: sourceDir() + "/file1",
		},
		{
			Path:    homeDir() + "/dir1/nestedDir/file4",
			Content: sourceDir() + "/dir1/nestedDir/file4",
		},
		{
			Path:    homeDir() + "/someFileInstalledInAPreviousRun",
			Content: sourceDir() + "/file1",
		},
	}
	dootCache.Save()

	install.Install()
	assertHomeDirContents(t, "", []string{
		"file2.txt",
		"dir1",
		"dir3",
	})
	assertHomeDirContents(t, "dir1", []string{
		"file3",
	})

	homePath := NewAbsolutePath(homeDir())
	assertCache(t, []AssertCacheEntry{
		{Path: homePath.Join("file1"), Content: sourceDir() + "/file1"},
		{Path: homePath.Join("file2.txt"), Content: sourceDir() + "/file2.txt"},
		{Path: homePath.Join("dir1/file3"), Content: sourceDir() + "/dir1/file3"},
		{Path: homePath.Join("dir1/nestedDir/file4"), Content: sourceDir() + "/dir1/nestedDir/file4"},
		{Path: homePath.Join("dir3/file6"), Content: sourceDir() + "/dir3/file6"},
	})
}

func TestInstall_IncrementalUpdateLink(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestInstall(t, config)
	createSymlink(homeDir(), "file1", sourceDir()+"/correct-target")
	createSymlink(homeDir(), "file2.txt", sourceDir()+"/incorrect-target")

	dootCache := cache.Load()
	cacheEntry := dootCache.GetEntry(sourceDir() + ":" + homeDir())
	cacheEntry.Links = []*cache.InstalledFile{
		{
			Path:    homeDir() + "/file1",
			Content: sourceDir() + "/file1", // Lie to the cache and see that the link is not updated
		},
		{
			Path:    homeDir() + "/file2.txt",
			Content: sourceDir() + "/incorrect-target",
		},
	}
	dootCache.Save()

	install.Install()
	assertHomeLink(t, "file1", sourceDir()+"/correct-target") // Not updated
	assertHomeLink(t, "file2.txt", sourceDir()+"/file2.txt")
}

func TestInstall_SilentOverwrite(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestInstall(t, config)
	// Due to the implementation of File(), it has the same contents as the one in dotfiles dir
	createFile(homeDir(), File("file1"))
	createSymlink(homeDir(), "file2.txt", sourceDir()+"/file2.txt")
	assertHomeRegularFile(t, "file1")

	install.Install()
	assertHomeLink(t, "file1", sourceDir()+"/file1")
	assertHomeDirContents(t, "", []string{
		"file1",
		"file2.txt",
		"dir1",
		"dir3",
	})

	homePath := NewAbsolutePath(homeDir())
	assertCache(t, []AssertCacheEntry{
		{Path: homePath.Join("file1"), Content: sourceDir() + "/file1"},
		{Path: homePath.Join("file2.txt"), Content: sourceDir() + "/file2.txt"},
		{Path: homePath.Join("dir1/file3"), Content: sourceDir() + "/dir1/file3"},
		{Path: homePath.Join("dir1/nestedDir/file4"), Content: sourceDir() + "/dir1/nestedDir/file4"},
		{Path: homePath.Join("dir3/file6"), Content: sourceDir() + "/dir3/file6"},
	})
}

func TestInstall_OverwriteN(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestInstall(t, config)
	createFile(homeDir(), FsFile{Name: "file1", Content: "This is an outdated text"})
	createSymlink(homeDir(), "file2.txt", sourceDir()+"/outdatedLink")

	response := "n"
	utils.USER_INPUT_MOCK_RESPONSE = &response

	install.Install()
	assertHomeRegularFile(t, "file1")
	assertHomeLink(t, "file2.txt", sourceDir()+"/outdatedLink")

	homePath := NewAbsolutePath(homeDir())
	assertCache(t, []AssertCacheEntry{
		{Path: homePath.Join("dir1/file3"), Content: sourceDir() + "/dir1/file3"},
		{Path: homePath.Join("dir1/nestedDir/file4"), Content: sourceDir() + "/dir1/nestedDir/file4"},
		{Path: homePath.Join("dir3/file6"), Content: sourceDir() + "/dir3/file6"},
	})
}

func TestInstall_OverwriteY(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestInstall(t, config)
	createFile(homeDir(), FsFile{Name: "file1", Content: "This is an outdated text"})
	createSymlink(homeDir(), "file2.txt", sourceDir()+"/outdatedLink")

	response := "y"
	utils.USER_INPUT_MOCK_RESPONSE = &response

	install.Install()
	assertHomeLink(t, "file1", sourceDir()+"/file1")
	assertHomeLink(t, "file2.txt", sourceDir()+"/file2.txt")

	homePath := NewAbsolutePath(homeDir())
	assertCache(t, []AssertCacheEntry{
		{Path: homePath.Join("file1"), Content: sourceDir() + "/file1"},
		{Path: homePath.Join("file2.txt"), Content: sourceDir() + "/file2.txt"},
		{Path: homePath.Join("dir1/file3"), Content: sourceDir() + "/dir1/file3"},
		{Path: homePath.Join("dir1/nestedDir/file4"), Content: sourceDir() + "/dir1/nestedDir/file4"},
		{Path: homePath.Join("dir3/file6"), Content: sourceDir() + "/dir3/file6"},
	})
}

func TestInstall_WithCryptInitialized(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestInstall(t, config)
	initializeGitCrypt()

	install.Install()
	assertHomeDirContents(t, "dir3", []string{
		"file6",
		"file7",
	})
	assertHomeLink(t, "dir3/file7", sourceDir()+"/dir3/file7.doot-crypt")

	install.Clean()
	assertHomeDirContents(t, "", []string{})
}

func TestInstall_WithHostSpecificDir(t *testing.T) {
	myHost, err := os.Hostname()
	assert.NoError(t, err)

	config := config.DefaultConfig()
	config.ExcludeFiles = []string{".git"}
	config.Hosts = map[string]string{
		"other_host1": "hosts/OTHER",
		myHost:        "hosts/HOST",
	}
	setUpFiles_TestInstall(t, config)
	initializeGitCrypt()
	createNode(sourceDir(), Dir("hosts", []FsNode{
		Dir("OTHER", []FsNode{
			File("inOtherHost"),
			Dir("inOtherHostDir", []FsNode{File("inOtherHostDirFile")}),
		}),
		Dir("HOST", []FsNode{
			File("inMyHost"),
			Dir("inMyHostDir", []FsNode{File("inMyHostDirFile")}),
			File("file2.doot-crypt.txt"),
			Dir("dir2", []FsNode{File("file5")}),
		}),
	}))

	install.Install()
	assertHomeDirContents(t, "", []string{
		".file1",
		".file2.txt",
		".dir1",
		".dir2",
		".dir3",
		".inMyHost",
		".inMyHostDir",
	})
	assertHomeDirContents(t, ".inMyHostDir", []string{"inMyHostDirFile"})
	assertHomeLink(t, ".inMyHost", sourceDir()+"/hosts/HOST/inMyHost")
	assertHomeLink(t, ".file2.txt", sourceDir()+"/hosts/HOST/file2.doot-crypt.txt")
	assertHomeLink(t, ".dir1/file3", sourceDir()+"/dir1/file3")
	assertHomeLink(t, ".dir2/file5", sourceDir()+"/hosts/HOST/dir2/file5")

	install.Clean()
	assertHomeDirContents(t, "", []string{})
}

func TestInstall_DoNotRemoveUnexpectedFiles(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestInstall(t, config)

	install.Install()
	assertHomeDirContents(t, "", []string{
		"file1",
		"file2.txt",
		"dir1",
		"dir3",
	})
	// User manually changes some files, which should NOT be removed when cleaning
	os.Remove(homeDir() + "/file1")
	createNode(homeDir(), File("file1"))
	os.Remove(homeDir() + "/file2.txt")
	createSymlink(homeDir(), "file2.txt", homeDir()+"/incorrect_link") // Link does not point to dotfiles dir

	install.Clean()
	assertHomeDirContents(t, "", []string{
		"file1",
		"file2.txt",
	})
}

func TestInstall_ExploreExcludedDirs(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	config.ExploreExcludedDirs = true
	config.ExcludeFiles = []string{"dir*", "**/nestedDir", "**/.*"}
	config.IncludeFiles = []string{"dir1/nestedDir/file4", "dir3/file6", ".dir2/file5"}
	setUpFiles_TestInstall(t, config)

	install.Install()
	assertHomeDirContents(t, "", []string{
		"file1",
		"file2.txt",
		"dir1",
		".dir2",
		"dir3",
	})
	assertHomeDirContents(t, "dir1", []string{
		"nestedDir",
	})
	assertHomeDirContents(t, "dir1/nestedDir", []string{
		"file4",
	})
	assertHomeDirContents(t, ".dir2", []string{
		"file5",
	})
	assertHomeDirContents(t, "dir3", []string{
		"file6",
	})
}

func setUpFiles_TestInstall(t *testing.T, config config.Config) {
	SetUpFiles(t, []FsNode{
		Dir("doot", []FsNode{
			ConfigFile(config),
		}),
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
		}),
		Dir("dir3", []FsNode{
			File("file6"),
			File("file7.doot-crypt"),
		}),
		Dir("emptyDir", []FsNode{}),
	})
}
