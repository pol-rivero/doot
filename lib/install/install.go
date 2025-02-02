package install

import (
	"github.com/pol-rivero/doot/lib/cache"
	"github.com/pol-rivero/doot/lib/config"
)

func Install() {
	dotfilesDir := FindDotfilesDir()
	config := config.FromDotfilesDir(dotfilesDir)
	cache := cache.Load()
	installedFilesCache := cache.UseDir(dotfilesDir)

	ignoreDootCrypt := !gitCryptIsInitialized()
	filter := CreateFilter(&config, ignoreDootCrypt)
	fileList := ScanDirectory(dotfilesDir, filter)
	fileMapping := NewFileMapping(&config, fileList)

	fileMapping.InstallNewLinks(installedFilesCache.Targets)
	fileMapping.RemoveStaleLinks(installedFilesCache.Targets)

	installedFilesCache.Targets = fileMapping.GetTargets()
	cache.Save()
}

func gitCryptIsInitialized() bool {
	// TODO
	return false
}
