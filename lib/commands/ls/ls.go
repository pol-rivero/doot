package ls

import (
	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/cache"
	"github.com/pol-rivero/doot/lib/common/config"
	"github.com/pol-rivero/doot/lib/common/log"
)

func ListInstalledFiles(asJson bool) {
	dotfilesDir := common.FindDotfilesDir()
	config := config.FromDotfilesDir(dotfilesDir)

	cacheKey := cache.ComputeCacheKey(dotfilesDir, config.TargetDir)
	cache := cache.Load()
	installedFilesCache := cache.GetEntry(cacheKey)

	installedLinks := installedFilesCache.GetLinks()
	if asJson {
		log.Printlnf("%s", installedLinks.ToJson())
	} else {
		log.Printlnf("%s", installedLinks.PrintList())
	}
}
