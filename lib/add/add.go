package add

import (
	"os"
	"path/filepath"

	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/config"
	"github.com/pol-rivero/doot/lib/common/glob_collection"
	"github.com/pol-rivero/doot/lib/common/log"
	"github.com/pol-rivero/doot/lib/crypt"
	"github.com/pol-rivero/doot/lib/install"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/pol-rivero/doot/lib/utils/set"
)

func Add(files []string, isCrypt bool, isHostSpecific bool) {
	dotfilesDir := common.FindDotfilesDir()
	config := config.FromDotfilesDir(dotfilesDir)
	params := ProcessAddedFileParams{
		crypt:             isCrypt,
		hostSpecificDir:   getHostSpecificDir(&config, isHostSpecific),
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

	for _, file := range files {
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
		// Hardlink instead of copy, the original file will be replaced on install anyway
		err = os.Link(file, dotfilePath.Str())
		if err == nil {
			log.Info("Created hardlink %s -> %s", file, dotfilePath)
		} else if os.IsExist(err) {
			handleDotfileAlreadyExists(file, dotfilePath)
		} else {
			log.Error("Error hardlinking %s to %s: %v", file, dotfilePath, err)
		}
	}

	log.Info("Files have been copied to the dotfiles directory, now running 'install'...")
	install.Install(false)
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
%s = "%s-files"`, hostname, hostname, hostname)
	}
	return hostSpecificDir
}

func handleDotfileAlreadyExists(targetFile string, dotfilePath AbsolutePath) {
	fileInfo, err := os.Lstat(targetFile)
	if err != nil {
		panic("Should have been checked in ProcessAddedFile")
	}
	if fileInfo.Mode()&os.ModeSymlink != 0 && getSymlinkTarget(targetFile) == dotfilePath.Str() {
		log.Warning("Link %s -> %s already exists, skipping", targetFile, dotfilePath)
	} else {
		log.Error("Dotfile %s already exists. If you really want to overwrite it, delete it first", dotfilePath)
	}
}

func getSymlinkTarget(linkPath string) string {
	linkSource, linkErr := os.Readlink(linkPath)
	if linkErr != nil {
		log.Fatal("Failed to read link %s: %v", linkPath, linkErr)
	}
	return linkSource
}
