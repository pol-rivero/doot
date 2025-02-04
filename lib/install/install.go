package install

import (
	"path/filepath"

	"github.com/pol-rivero/doot/lib/cache"
	"github.com/pol-rivero/doot/lib/config"
)

func Install() {
	dotfilesDir := FindDotfilesDir()
	config := config.FromDotfilesDir(dotfilesDir)

	cacheKey := dotfilesDir.Str() + string(filepath.ListSeparator) + config.TargetDir
	cache := cache.Load()
	installedFilesCache := cache.GetEntry(cacheKey)

	ignoreDootCrypt := !gitCryptIsInitialized()
	filter := CreateFilter(&config, ignoreDootCrypt)
	fileList := ScanDirectory(dotfilesDir, filter)
	fileMapping := NewFileMapping(dotfilesDir, &config, fileList)

	fileMapping.InstallNewLinks(installedFilesCache.GetTargets())
	fileMapping.RemoveStaleLinks(installedFilesCache.GetTargets())

	installedFilesCache.SetTargets(fileMapping.GetInstalledTargets())
	cache.Save()
}

func gitCryptIsInitialized() bool {
	// TODO
	return false
}
