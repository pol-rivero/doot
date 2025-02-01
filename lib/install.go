package lib

import (
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/pol-rivero/doot/lib/cache"
	"github.com/pol-rivero/doot/lib/config"
	"github.com/pol-rivero/doot/lib/helpers"
)

func Install() {
	dotfilesDir := helpers.FindDotfilesDir()
	config := config.FromDotfilesDir(dotfilesDir)
	cache := cache.Load()
	installedFilesCache := cache.UseDir(dotfilesDir)

	ignoreDootCrypt := !gitCryptIsInitialized()
	filter := helpers.CreateFilter(&config, ignoreDootCrypt)
	fileList := helpers.ScanDirectory(dotfilesDir, filter)
	fileMapping := getFileMapping(fileList, &config)

	fileMapping.InstallNewLinks(installedFilesCache.Targets)
	fileMapping.RemoveStaleLinks(installedFilesCache.Targets)

	installedFilesCache.Targets = fileMapping.GetTargets()
	cache.Save()
}

func getFileMapping(sourceFiles []string, config *config.Config) helpers.FileMapping {
	mapping := helpers.NewFileMapping(len(sourceFiles))
	targetBaseDir := config.TargetDir

	for _, sourceFile := range sourceFiles {
		target := path.Join(targetBaseDir, mapSourceToTarget(sourceFile, config.ImplicitDot, helpers.NewSetFromSlice(config.ImplicitDotIgnore)))
		mapping.Add(sourceFile, target)
	}
	return mapping
}

func mapSourceToTarget(relativePath string, implicitDot bool, implicitDotIgnore helpers.Set[string]) string {
	target := relativePath
	if implicitDot && !implicitDotIgnore.Contains(getTopLevelDir(relativePath)) {
		target = "." + relativePath
	}
	log.Printf("Mapping %s to %s", relativePath, target)
	return target
}

func getTopLevelDir(filePath string) string {
	firstSeparatorIndex := strings.IndexRune(filePath, filepath.Separator)
	if firstSeparatorIndex == -1 {
		return filePath
	}
	return filePath[:firstSeparatorIndex]
}

func gitCryptIsInitialized() bool {
	// TODO
	return false
}
