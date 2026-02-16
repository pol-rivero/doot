package add

import (
	"os"
	"path/filepath"

	"github.com/pol-rivero/doot/lib/commands/crypt"
	"github.com/pol-rivero/doot/lib/commands/install"
	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/cache"
	"github.com/pol-rivero/doot/lib/common/config"
	"github.com/pol-rivero/doot/lib/common/glob_collection"
	"github.com/pol-rivero/doot/lib/common/log"
	"github.com/pol-rivero/doot/lib/linkmode"
	. "github.com/pol-rivero/doot/lib/types"
	file_utils "github.com/pol-rivero/doot/lib/utils/files"
	"github.com/pol-rivero/doot/lib/utils/set"
)

func Add(files []string, isCrypt bool, isHostSpecific bool) {
	dotfilesDir := common.FindDotfilesDir()
	config := config.FromDotfilesDir(dotfilesDir)
	linkMode := linkmode.GetLinkMode(&config)

	cacheKey := cache.ComputeCacheKey(dotfilesDir, config.TargetDir)
	cache := cache.Load()
	installedLinks := cache.GetEntry(cacheKey).GetLinks()

	params := ProcessAddedFileParams{
		crypt:             isCrypt,
		hostSpecificDir:   getHostSpecificDir(&config, isHostSpecific),
		dotfilesDir:       dotfilesDir.Str(),
		targetDir:         config.TargetDir,
		implicitDot:       config.ImplicitDot,
		implicitDotIgnore: set.NewFromSlice(config.ImplicitDotIgnore),
		includeFiles:      glob_collection.NewGlobCollection(config.IncludeFiles),
		excludeFiles:      glob_collection.NewGlobCollection(config.ExcludeFiles),
	}

	if isCrypt && !crypt.GitCryptIsInitialized(dotfilesDir) {
		log.Error("Can't add private files with --crypt flag because repository is not initialized. Run 'doot crypt init' first.")
		return
	}

	addedFiles := make([]AbsolutePath, 0, 16)
	for _, file := range files {
		if alreadyManaged(file, &installedLinks, linkMode) {
			log.Printlnf("%s is already managed by doot", file)
			continue
		}
		dotfileRelativePath, err := ProcessAddedFile(file, params)
		if err != nil {
			log.Error("Can't add %s: %v", file, err)
			continue
		}
		dotfilePath := dotfilesDir.JoinPath(dotfileRelativePath)
		err = os.MkdirAll(filepath.Dir(dotfilePath.Str()), 0755)
		if err != nil {
			log.Error("Error creating directory %s: %v", filepath.Dir(dotfilePath.Str()), err)
			continue
		}
		// Prefer hardlinking to avoid copying large files. The original file will be replaced on install anyway.
		// The dotfiles directory may reside in a different filesystem than the home directory, fallback to copy if hardlinking fails.
		err = file_utils.HardlinkOrCopyFile(file, dotfilePath.Str(), false)
		if err == nil {
			log.Info("Copied file %s -> %s", file, dotfilePath)
			addedFiles = append(addedFiles, RelativeToPWD(file))
		} else if os.IsExist(err) {
			log.Error("Dotfile %s already exists. If you really want to overwrite it, delete it first", dotfilePath)
		} else {
			log.Error("Error copying %s to %s: %v", file, dotfilePath, err)
		}
	}

	log.Info("Files have been copied to the dotfiles directory, now running 'install'...")
	install.InstallAfterAdd(false, addedFiles)
}

func getHostSpecificDir(config *config.Config, isHostSpecific bool) string {
	if !isHostSpecific {
		return ""
	}
	hostname, err := os.Hostname()
	if err != nil {
		log.Error("Error getting hostname: %v", err)
		return ""
	}
	hostSpecificDir, ok := config.Hosts[hostname]
	if !ok {
		log.Fatal(`--host flag is set but your hostname (%s) is not in the hosts map. Consider adding the following to your doot config:
[hosts]
"%s" = "%s-files"`, hostname, hostname, hostname)
	}
	return hostSpecificDir
}

func alreadyManaged(file string, installedLinks *SymlinkCollection, linkMode linkmode.LinkMode) bool {
	installedLink := RelativeToPWD(file)
	dotfilePath := installedLinks.Get(installedLink)
	if dotfilePath.IsEmpty() {
		return false
	}
	return linkMode.IsInstalledLinkOf(installedLink.Str(), dotfilePath.Value())
}
