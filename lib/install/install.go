package install

import (
	"path/filepath"

	"github.com/fatih/color"
	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/cache"
	"github.com/pol-rivero/doot/lib/common/config"
	"github.com/pol-rivero/doot/lib/common/log"
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

	fileList := getFiles(&config, dotfilesDir)
	fileMapping := NewFileMapping(dotfilesDir, &config, fileList)

	added := fileMapping.InstallNewLinks(installedFilesCache.GetTargets())
	removed := fileMapping.RemoveStaleLinks(installedFilesCache.GetTargets())

	installedFilesCache.SetTargets(fileMapping.GetInstalledTargets())
	cache.Save()

	printChanges(added, removed)
}

func printChanges(added int, removed int) {
	if added == 0 && removed == 0 {
		log.Printlnf("No changes made")
		return
	}
	if added > 0 {
		boldGreen := color.New(color.FgGreen, color.Bold).SprintFunc()
		log.Printf(boldGreen("%d")+color.GreenString(" %s added"), added, links(added))
	}
	if added > 0 && removed > 0 {
		log.Printf(", ")
	}
	if removed > 0 {
		boldRed := color.New(color.FgRed, color.Bold).SprintFunc()
		log.Printf(boldRed("%d")+color.RedString(" %s removed"), removed, links(removed))
	}
	log.Printlnf("")
}

func links(num int) string {
	if num == 1 {
		return "link"
	}
	return "links"
}
