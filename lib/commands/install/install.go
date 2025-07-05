package install

import (
	"path/filepath"

	"github.com/pol-rivero/doot/lib/commands/crypt"
	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/cache"
	"github.com/pol-rivero/doot/lib/common/config"
	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
)

type GetFilesFunc func(*config.Config, AbsolutePath) []RelativePath

func Install(fullClean bool) {
	getFiles := func(config *config.Config, dotfilesDir AbsolutePath) []RelativePath {
		ignoreDootCrypt := !crypt.GitCryptIsInitialized(dotfilesDir)
		filter := CreateFilter(config, ignoreDootCrypt)
		return ScanDirectory(dotfilesDir, &filter)
	}
	installImpl(getFiles, fullClean)
}

func Clean(fullClean bool) {
	getFiles := func(config *config.Config, dotfilesDir AbsolutePath) []RelativePath {
		return []RelativePath{}
	}
	installImpl(getFiles, fullClean)
}

func installImpl(getFiles GetFilesFunc, fullClean bool) {
	dotfilesDir := common.FindDotfilesDir()
	config := config.FromDotfilesDir(dotfilesDir)

	if config.UseHardlinks && fullClean {
		log.Fatal("--full-clean is not supported because you have set use_hardlinks=true")
	}

	cacheKey := dotfilesDir.Str() + string(filepath.ListSeparator) + config.TargetDir
	cache := cache.Load()
	installedFilesCache := cache.GetEntry(cacheKey)
	if fullClean {
		recalculateCache(installedFilesCache, dotfilesDir, config.TargetDir)
	}

	common.RunHooks(dotfilesDir, "before-update")
	fileList := getFiles(&config, dotfilesDir)
	fileMapping := NewFileMapping(dotfilesDir, &config, fileList)

	oldLinks := installedFilesCache.GetLinks()
	added := fileMapping.InstallNewLinks()
	removed := fileMapping.RemoveStaleLinks(&oldLinks)

	installedFilesCache.SetLinks(fileMapping.GetInstalledTargets())
	cache.Save()

	common.RunHooks(dotfilesDir, "after-update")
	printChanges(added, removed)
}
