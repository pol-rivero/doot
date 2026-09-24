package restore

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/cache"
	"github.com/pol-rivero/doot/lib/common/config"
	"github.com/pol-rivero/doot/lib/common/log"
	"github.com/pol-rivero/doot/lib/linkmode"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/pol-rivero/doot/lib/utils/files"
)

func Restore(inputFiles []string) {
	dotfilesDir := common.FindDotfilesDir()
	config := config.FromDotfilesDir(dotfilesDir)

	cacheKey := cache.ComputeCacheKey(dotfilesDir, config.TargetDir)
	cache := cache.Load()
	installedFilesCache := cache.GetEntry(cacheKey)

	installedLinks := installedFilesCache.GetLinks()
	linkMode := linkmode.GetLinkMode(&config)
	successCount := restoreFiles(inputFiles, installedLinks, dotfilesDir, linkMode)

	installedFilesCache.SetLinks(installedLinks)
	cache.Save()

	if successCount == 0 {
		os.Exit(1)
	} else {
		fileOrFiles := map[bool]string{true: "files", false: "file"}[successCount != 1]
		log.Printlnf("Successfully restored %d %s", successCount, fileOrFiles)
	}
}

func restoreFiles(inputFiles []string, installedLinks SymlinkCollection, dotfilesDir AbsolutePath, linkMode linkmode.LinkMode) int {
	successCount := 0
	for _, rawInput := range inputFiles {
		filePath, err := ensureFileExists(rawInput)
		if err == nil {
			err = restoreFile(filePath, installedLinks, dotfilesDir, linkMode)
		}

		if err != nil {
			log.Error("Failed to restore '%s': %v", rawInput, err)
		} else {
			log.Info("Successfully restored '%s'", rawInput)
			successCount++
		}
	}
	return successCount
}

func ensureFileExists(rawInput string) (AbsolutePath, error) {
	cleanAbsFile, err := filepath.Abs(rawInput)
	if err != nil {
		log.Fatal("Failed to get absolute path for '%s': %v", rawInput, err)
	}
	info, err := os.Lstat(cleanAbsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", errors.New("file not found")
		}
		return "", err
	}
	if info.IsDir() {
		return "", errors.New("it's a directory, you must specify files")
	}
	return NewAbsolutePath(cleanAbsFile), nil
}

func restoreFile(filePath AbsolutePath, installedLinks SymlinkCollection, dotfilesDir AbsolutePath, linkMode linkmode.LinkMode) error {
	for linkPath, linkContent := range installedLinks.Iter() {
		if linkPath == filePath || linkContent == filePath {
			err := overwriteLink(linkPath, linkContent, dotfilesDir, linkMode)
			if err == nil {
				installedLinks.Remove(linkPath)
			}
			return err
		}
	}
	return errors.New("it's not a dotfile managed by doot")
}

func overwriteLink(symlinkPath, dotfilePath, dotfilesDir AbsolutePath, linkMode linkmode.LinkMode) error {
	if err := ensureLinkCanBeOverwritten(symlinkPath, dotfilePath, linkMode); err != nil {
		return err
	}
	log.Info("Moving '%s' -> '%s'", dotfilePath, symlinkPath)
	if err := files.MoveOrCopyFile(dotfilePath.Str(), symlinkPath.Str(), true); err != nil {
		return err
	}
	files.CleanupEmptyDir(dotfilePath.Parent(), dotfilesDir)
	return nil
}

func ensureLinkCanBeOverwritten(symlinkPath, dotfilePath AbsolutePath, linkMode linkmode.LinkMode) error {
	_, err := os.Lstat(symlinkPath.Str())
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}
	if !linkMode.IsInstalledLinkOf(symlinkPath.Str(), dotfilePath) {
		return fmt.Errorf("'%s' was modified externally and is no longer a link to '%s'. Restoring it would overwrite its changes, please review it manually", symlinkPath, dotfilePath)
	}
	return nil
}
