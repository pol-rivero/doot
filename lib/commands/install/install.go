package install

import (
	"github.com/pol-rivero/doot/lib/commands/crypt"
	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/cache"
	"github.com/pol-rivero/doot/lib/common/config"
	"github.com/pol-rivero/doot/lib/linkmode"
	. "github.com/pol-rivero/doot/lib/types"
)

type GetFilesFunc func(*config.Config, AbsolutePath) []RelativePath

func Install(fullClean bool) {
	regularInstall(fullClean, nil)
}

func InstallAfterAdd(fullClean bool, extraAddedFiles []AbsolutePath) {
	regularInstall(fullClean, extraAddedFiles)
}

func regularInstall(fullClean bool, extraAddedFiles []AbsolutePath) {
	getFiles := func(config *config.Config, dotfilesDir AbsolutePath) []RelativePath {
		ignoreDootCrypt := !crypt.GitCryptIsInitialized(dotfilesDir)
		filter := CreateFilter(config, ignoreDootCrypt)
		return ScanDirectory(dotfilesDir, &filter)
	}
	added, removed := install(getFiles, fullClean)
	printChanges(added, removed, extraAddedFiles)
}

func Clean(fullClean bool) {
	getFiles := func(config *config.Config, dotfilesDir AbsolutePath) []RelativePath {
		return []RelativePath{}
	}
	added, removed := install(getFiles, fullClean)
	printChanges(added, removed, nil)
}

func install(getFiles GetFilesFunc, fullClean bool) ([]AbsolutePath, []AbsolutePath) {
	dotfilesDir := common.FindDotfilesDir()
	config := config.FromDotfilesDir(dotfilesDir)
	linkMode := linkmode.GetLinkMode(&config)

	cacheKey := cache.ComputeCacheKey(dotfilesDir, config.TargetDir)
	cache := cache.Load()
	installedFilesCache := cache.GetEntry(cacheKey)
	if fullClean {
		installedFilesCache.Links = linkMode.RecalculateCache(dotfilesDir, config.TargetDir)
	}

	common.RunHooks(dotfilesDir, "before-update")
	fileList := getFiles(&config, dotfilesDir)
	fileMapping := NewFileMapping(dotfilesDir, &config, fileList)

	oldLinks := installedFilesCache.GetLinks()
	removed := fileMapping.RemoveStaleLinks(&oldLinks)
	added := fileMapping.InstallNewLinks()

	installedFilesCache.SetLinks(fileMapping.GetInstalledTargets())
	cache.Save()

	common.RunHooks(dotfilesDir, "after-update")
	return added, removed
}
