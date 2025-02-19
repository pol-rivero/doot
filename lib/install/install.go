package install

import (
	"path/filepath"

	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/cache"
	"github.com/pol-rivero/doot/lib/common/config"
	"github.com/pol-rivero/doot/lib/crypt"
	. "github.com/pol-rivero/doot/lib/types"
)

type GetFilesFunc func(*config.Config, AbsolutePath) []RelativePath

func Install() {
	getFiles := func(config *config.Config, dotfilesDir AbsolutePath) []RelativePath {
		ignoreDootCrypt := !crypt.GitCryptIsInitialized(dotfilesDir)
		filter := CreateFilter(config, ignoreDootCrypt)
		return ScanDirectory(dotfilesDir, &filter)
	}
	installImpl(getFiles)
}

func Clean() {
	getFiles := func(config *config.Config, dotfilesDir AbsolutePath) []RelativePath {
		return []RelativePath{}
	}
	installImpl(getFiles)
}

func installImpl(getFiles GetFilesFunc) {
	dotfilesDir := common.FindDotfilesDir()
	config := config.FromDotfilesDir(dotfilesDir)

	cacheKey := dotfilesDir.Str() + string(filepath.ListSeparator) + config.TargetDir
	cache := cache.Load()
	installedFilesCache := cache.GetEntry(cacheKey)

	common.RunHooks(dotfilesDir, "before-update")
	fileList := getFiles(&config, dotfilesDir)
	fileMapping := NewFileMapping(dotfilesDir, &config, fileList)

	oldLinks := installedFilesCache.GetLinks()
	added := fileMapping.InstallNewLinks(&oldLinks)
	removed := fileMapping.RemoveStaleLinks(&oldLinks)

	installedFilesCache.SetLinks(fileMapping.GetInstalledTargets())
	cache.Save()

	common.RunHooks(dotfilesDir, "after-update")
	printChanges(added, removed)
}
