package test

import (
	"os"
	"testing"

	"github.com/pol-rivero/doot/lib/commands/install"
	"github.com/pol-rivero/doot/lib/common/cache"
	"github.com/pol-rivero/doot/lib/common/config"
	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/pol-rivero/doot/lib/utils"
	"github.com/stretchr/testify/assert"
)

func TestInstall_DefaultConfig(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestInstall(t, config)

	install.Install(false)
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
	assertHomeSymlink(t, "file1", sourceDir()+"/file1")
	assertHomeSymlink(t, "dir1/nestedDir/file4", sourceDir()+"/dir1/nestedDir/file4")

	os.Remove(sourceDir() + "/file1")
	install.Install(false)
	assertHomeDirContents(t, "", []string{
		"file2.txt",
		"dir1",
		"dir3",
	})

	install.Clean(false)
	assertHomeDirContents(t, "", []string{})
}

func TestInstall_Hardlinks(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	config.UseHardlinks = true
	setUpFiles_TestInstall(t, config)

	install.Install(false)
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
	assertHomeHardlink(t, "file1", sourceDir()+"/file1")
	assertHomeHardlink(t, "dir1/nestedDir/file4", sourceDir()+"/dir1/nestedDir/file4")

	os.Remove(sourceDir() + "/file1")
	install.Install(false)
	assertHomeDirContents(t, "", []string{
		"file2.txt",
		"dir1",
		"dir3",
	})

	install.Clean(false)
	assertHomeDirContents(t, "", []string{})
}

func TestInstall_HiddenFiles(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	config.ExcludeFiles = []string{"file2.txt"}
	setUpFiles_TestInstall(t, config)

	install.Install(false)
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
	assertHomeSymlink(t, ".dir2/file5", sourceDir()+"/.dir2/file5")

	install.Clean(false)
	assertHomeDirContents(t, "", []string{})
}

func TestInstall_ImplicitDot(t *testing.T) {
	config := config.DefaultConfig()
	config.ExcludeFiles = []string{}
	config.ImplicitDot = true
	config.ImplicitDotIgnore = []string{"file2.txt", "dir3"}
	setUpFiles_TestInstall(t, config)

	install.Install(false)
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
	assertHomeSymlink(t, ".file1", sourceDir()+"/file1")
	assertHomeSymlink(t, ".dir1/file3", sourceDir()+"/dir1/file3")

	install.Clean(false)
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

	install.Install(false)
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

	install.Clean(false)
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
	install.Install(false)

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

	install.Clean(false)
	assertCache(t, []AssertCacheEntry{})
}

func TestInstall_IncrementalInstall(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestInstall(t, config)
	createSymlink(homeDir(), "someFileInstalledInAPreviousRun", sourceDir()+"/file1")

	// someFileInstalledInAPreviousRun is no longer in dotfiles dir, so it should be removed
	// file1 and file4 are already installed, but the disk is checked to catch changes like these
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

	install.Install(false)
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
	createSymlink(homeDir(), "file1", "/incorrect-target")

	dootCache := cache.Load()
	cacheEntry := dootCache.GetEntry(sourceDir() + ":" + homeDir())
	cacheEntry.Links = []*cache.InstalledFile{
		{
			Path:    homeDir() + "/file1",
			Content: sourceDir() + "/file1", // Even though it's in the cache, it should be checked for changes
		},
	}
	dootCache.Save()

	utils.USER_INPUT_MOCK_RESPONSE = "n"
	install.Install(false)
	assertHomeSymlink(t, "file1", "/incorrect-target")

	utils.USER_INPUT_MOCK_RESPONSE = "y"
	install.Install(false)
	assertHomeSymlink(t, "file1", sourceDir()+"/file1")
}

func TestInstall_IncrementalUpdateSymlinkToHardlink(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	config.UseHardlinks = true
	setUpFiles_TestInstall(t, config)
	createSymlink(homeDir(), "file1", "/incorrect-target")

	dootCache := cache.Load()
	cacheEntry := dootCache.GetEntry(sourceDir() + ":" + homeDir())
	cacheEntry.Links = []*cache.InstalledFile{
		{
			Path:    homeDir() + "/file1",
			Content: sourceDir() + "/file1",
		},
	}
	dootCache.Save()

	utils.USER_INPUT_MOCK_RESPONSE = "n"
	install.Install(false)
	assertHomeSymlink(t, "file1", "/incorrect-target")

	utils.USER_INPUT_MOCK_RESPONSE = "y"
	install.Install(false)
	assertHomeHardlink(t, "file1", sourceDir()+"/file1")
}

func TestInstall_SilentOverwrite(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestInstall(t, config)
	// Due to the implementation of File(), it has the same contents as the one in dotfiles dir
	createFile(homeDir(), File("file1"))
	createSymlink(homeDir(), "file2.txt", sourceDir()+"/file2.txt")
	assertHomeRegularFile(t, "file1")

	install.Install(false)
	assertHomeSymlink(t, "file1", sourceDir()+"/file1")
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

func TestInstall_SilentOverwriteSymlinkToHardlink(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	config.UseHardlinks = true
	setUpFiles_TestInstall(t, config)
	correctTarget := sourceDir() + "/file1"
	createSymlink(homeDir(), "file1", correctTarget)

	dootCache := cache.Load()
	cacheEntry := dootCache.GetEntry(sourceDir() + ":" + homeDir())
	cacheEntry.Links = []*cache.InstalledFile{
		{
			Path:    homeDir() + "/file1",
			Content: correctTarget,
		},
	}
	dootCache.Save()

	utils.USER_INPUT_MOCK_RESPONSE = "n"
	install.Install(false)
	assertHomeHardlink(t, "file1", correctTarget)
}

func TestInstall_SilentOverwriteHardlinkToSymlink(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestInstall(t, config)
	correctTarget := sourceDir() + "/file1"
	createHardlink(homeDir(), "file1", correctTarget)

	dootCache := cache.Load()
	cacheEntry := dootCache.GetEntry(sourceDir() + ":" + homeDir())
	cacheEntry.Links = []*cache.InstalledFile{
		{
			Path:    homeDir() + "/file1",
			Content: correctTarget,
		},
	}
	dootCache.Save()

	utils.USER_INPUT_MOCK_RESPONSE = "n"
	install.Install(false)
	assertHomeSymlink(t, "file1", correctTarget)
}

func TestInstall_OverwriteN(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestInstall(t, config)

	createFile(homeDir(), FsFile{Name: "file1", Content: "This is an outdated text"})
	createSymlink(homeDir(), "file2.txt", "/outdatedLink")

	utils.USER_INPUT_MOCK_RESPONSE = "n"
	install.Install(false)
	assertHomeRegularFile(t, "file1")
	assertHomeSymlink(t, "file2.txt", "/outdatedLink")

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
	createSymlink(homeDir(), "file2.txt", "/outdatedLink")

	utils.USER_INPUT_MOCK_RESPONSE = "y"
	install.Install(false)
	assertHomeSymlink(t, "file1", sourceDir()+"/file1")
	assertHomeSymlink(t, "file2.txt", sourceDir()+"/file2.txt")

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

	install.Install(false)
	assertHomeDirContents(t, "dir3", []string{
		"file6",
		"file7",
	})
	assertHomeSymlink(t, "dir3/file7", sourceDir()+"/dir3/file7.doot-crypt")

	install.Clean(false)
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

	install.Install(false)
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
	assertHomeSymlink(t, ".inMyHost", sourceDir()+"/hosts/HOST/inMyHost")
	assertHomeSymlink(t, ".file2.txt", sourceDir()+"/hosts/HOST/file2.doot-crypt.txt")
	assertHomeSymlink(t, ".dir1/file3", sourceDir()+"/dir1/file3")
	assertHomeSymlink(t, ".dir2/file5", sourceDir()+"/hosts/HOST/dir2/file5")

	install.Clean(false)
	assertHomeDirContents(t, "", []string{})
}

func TestInstall_DoNotRemoveUnexpectedFiles(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestInstall(t, config)

	install.Install(false)
	assertHomeDirContents(t, "", []string{
		"file1",
		"file2.txt",
		"dir1",
		"dir3",
	})
	// User manually changes some files, which should NOT be removed when cleaning
	os.Remove(homeDir() + "/file1")
	createNode(homeDir(), File("file1"))
	replaceWithSymlink(homeDir(), "file2.txt", homeDir()+"/incorrect_link") // Link does not point to dotfiles dir

	install.Clean(false)
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

	install.Install(false)
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

func TestInstall_Hooks(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestInstall(t, config)
	createHookFile("before-update", "before1.sh", `#!/bin/bash
		echo "before1 $PWD" >> before.txt`)
	createHookFile("before-update", "before2.sh", `#!/bin/bash
		echo "before2" >> before.txt`)
	createHookFile("after-update", "after1.sh", `#!/bin/bash
		echo "after" >> before.txt && echo "after" >> after.txt`)

	install.Install(false)
	assertHomeDirContents(t, "", []string{
		"before.txt",
		// after.txt should not be linked because it was created after the install
		"file1",
		"file2.txt",
		"dir1",
		"dir3",
	})
	assertHomeSymlink(t, "before.txt", sourceDir()+"/before.txt")
	beforeContent := readFile(sourceDir() + "/before.txt")
	assert.Equal(t, "before1 "+sourceDir()+"\nbefore2\nafter\n", beforeContent)
	afterContent := readFile(sourceDir() + "/after.txt")
	assert.Equal(t, "after\n", afterContent, "after hook was not executed")
}

func TestInstall_HooksFail(t *testing.T) {
	log.PanicInsteadOfExit = true
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestInstall(t, config)
	createHookFile("before-update", "before1.sh", `#!/bin/bash
		echo "i will fail" >> hook.txt
		exit 1`)
	createHookFile("before-update", "before2.sh", `#!/bin/bash
		echo "this should not be executed" >> hook.txt`)
	createHookFile("after-update", "after.sh", `#!/bin/bash
		echo "i shouldn't be executed either" >> hook.txt`)

	assert.Panics(t, func() {
		install.Install(false)
	})
	assertHomeDirContents(t, "", []string{}) // Install process should have been aborted
	beforeContent := readFile(sourceDir() + "/hook.txt")
	assert.Equal(t, "i will fail\n", beforeContent)
}

func TestInstall_FullClean(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestInstall(t, config)
	createNode(homeDir(), Dir("nested", []FsNode{Dir("dir", []FsNode{})}))
	createSymlink(homeDir()+"/nested/dir", "outdatedLink", sourceDir()+"/im-not-in-cache")

	install.Install(false)
	assert.FileExists(t, homeDir()+"/nested/dir/outdatedLink")

	install.Install(true)
	assert.NoFileExists(t, homeDir()+"/nested/dir/outdatedLink")
	assert.NoDirExists(t, homeDir()+"/nested")
}

func TestInstall_FullClean2(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestInstall(t, config)
	createNode(homeDir(), Dir("nested", []FsNode{Dir("dir", []FsNode{})}))
	createSymlink(homeDir()+"/nested/dir", "outdatedLink", sourceDir()+"/im-not-in-cache")

	install.Clean(false)
	assert.FileExists(t, homeDir()+"/nested/dir/outdatedLink")

	install.Clean(true)
	assertHomeDirContents(t, "", []string{})
}

func TestInstall_FullCleanHardlink(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	config.UseHardlinks = true
	setUpFiles_TestInstall(t, config)
	createNode(homeDir(), Dir("nested", []FsNode{Dir("dir", []FsNode{})}))
	createNode(homeDir(), File("doNotRemoveUnrelatedFile"))
	createHardlink(homeDir(), "doNotRemoveUnrelatedFile2", homeDir()+"/doNotRemoveUnrelatedFile")
	createHardlink(homeDir()+"/nested/dir", "outdatedLink", sourceDir()+"/file1")

	install.Install(false)
	assert.FileExists(t, homeDir()+"/nested/dir/outdatedLink")

	install.Install(true)
	assert.NoFileExists(t, homeDir()+"/nested/dir/outdatedLink")
	assert.NoDirExists(t, homeDir()+"/nested")
	assert.FileExists(t, homeDir()+"/doNotRemoveUnrelatedFile")
	assert.FileExists(t, homeDir()+"/doNotRemoveUnrelatedFile2")
}

func TestInstall_FullCleanHardlink2(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	config.UseHardlinks = true
	setUpFiles_TestInstall(t, config)
	createNode(homeDir(), Dir("nested", []FsNode{Dir("dir", []FsNode{})}))
	createHardlink(homeDir()+"/nested/dir", "outdatedLink", sourceDir()+"/file1")

	install.Clean(false)
	assert.FileExists(t, homeDir()+"/nested/dir/outdatedLink")

	install.Clean(true)
	assertHomeDirContents(t, "", []string{})
}

func TestInstall_AddDootCryptDoesntRequireConfirmation(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestInstall(t, config)
	initializeGitCrypt()

	install.Install(false)
	assertHomeSymlink(t, "file1", sourceDir()+"/file1")

	err := os.Rename(sourceDir()+"/file1", sourceDir()+"/file1.doot-crypt")
	assert.NoError(t, err)

	// Shouldn't wait for user input
	install.Install(false)
	assertHomeSymlink(t, "file1", sourceDir()+"/file1.doot-crypt")
}

func TestInstall_AdoptChanges(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestInstall(t, config)

	install.Install(false)
	assertHomeSymlink(t, "file1", sourceDir()+"/file1")
	assert.Equal(t, "dummy text for file file1", readFile(sourceDir()+"/file1"))

	os.Remove(homeDir() + "/file1")
	createFile(homeDir(), FsFile{Name: "file1", Content: "Some external program has replaced this"})
	assertHomeRegularFile(t, "file1")

	utils.USER_INPUT_MOCK_RESPONSE = "a"
	install.Install(false)
	assertHomeSymlink(t, "file1", sourceDir()+"/file1")
	assert.Equal(t, "Some external program has replaced this", readFile(sourceDir()+"/file1"))

	homePath := NewAbsolutePath(homeDir())
	assertCache(t, []AssertCacheEntry{
		{Path: homePath.Join("file1"), Content: sourceDir() + "/file1"},
		{Path: homePath.Join("file2.txt"), Content: sourceDir() + "/file2.txt"},
		{Path: homePath.Join("dir1/file3"), Content: sourceDir() + "/dir1/file3"},
		{Path: homePath.Join("dir1/nestedDir/file4"), Content: sourceDir() + "/dir1/nestedDir/file4"},
		{Path: homePath.Join("dir3/file6"), Content: sourceDir() + "/dir3/file6"},
	})
}

func TestInstall_AdoptChangesHardlink(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	config.UseHardlinks = true
	setUpFiles_TestInstall(t, config)

	install.Install(false)
	assertHomeHardlink(t, "file1", sourceDir()+"/file1")
	assert.Equal(t, "dummy text for file file1", readFile(sourceDir()+"/file1"))

	os.Remove(homeDir() + "/file1")
	createFile(homeDir(), FsFile{Name: "file1", Content: "Some external program has replaced this"})
	assertHomeRegularFile(t, "file1")

	utils.USER_INPUT_MOCK_RESPONSE = "a"
	install.Install(false)
	assertHomeHardlink(t, "file1", sourceDir()+"/file1")
	assert.Equal(t, "Some external program has replaced this", readFile(sourceDir()+"/file1"))

	homePath := NewAbsolutePath(homeDir())
	assertCache(t, []AssertCacheEntry{
		{Path: homePath.Join("file1"), Content: sourceDir() + "/file1"},
		{Path: homePath.Join("file2.txt"), Content: sourceDir() + "/file2.txt"},
		{Path: homePath.Join("dir1/file3"), Content: sourceDir() + "/dir1/file3"},
		{Path: homePath.Join("dir1/nestedDir/file4"), Content: sourceDir() + "/dir1/nestedDir/file4"},
		{Path: homePath.Join("dir3/file6"), Content: sourceDir() + "/dir3/file6"},
	})
}

func TestInstall_AdoptRegularFileReplacingSymlink(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false

	setUpFiles_TestInstall(t, config)
	replaceWithSymlink(sourceDir(), "file1", "/some-file")

	createFile(homeDir(), File("file1"))

	utils.USER_INPUT_MOCK_RESPONSE = "a"
	install.Install(false)
	assertHomeSymlink(t, "file1", sourceDir()+"/file1")
	assertRegularFile(t, sourceDir()+"/file1")
}

func TestInstall_AdoptRegularFileReplacingSymlinkUsingHardlinks(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	config.UseHardlinks = true

	setUpFiles_TestInstall(t, config)
	replaceWithSymlink(sourceDir(), "file1", "/some-file")

	createFile(homeDir(), File("file1"))

	utils.USER_INPUT_MOCK_RESPONSE = "a"
	install.Install(false)
	assertHomeHardlink(t, "file1", sourceDir()+"/file1")
	assertRegularFile(t, sourceDir()+"/file1")
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
