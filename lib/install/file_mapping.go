package install

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pol-rivero/doot/lib/config"
	"github.com/pol-rivero/doot/lib/log"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/pol-rivero/doot/lib/utils"
)

type FileMapping struct {
	mapping           map[AbsolutePath]AbsolutePath // Map target (symlink) -> source (dotfile)
	sourceBaseDir     AbsolutePath
	targetBaseDir     AbsolutePath
	implicitDot       bool
	implicitDotIgnore utils.Set[string]
}

func NewFileMapping(dotfilesDir AbsolutePath, config *config.Config, sourceFiles []RelativePath) FileMapping {
	mapping := FileMapping{
		mapping:           make(map[AbsolutePath]AbsolutePath, len(sourceFiles)),
		sourceBaseDir:     dotfilesDir,
		targetBaseDir:     NewAbsolutePath(config.TargetDir),
		implicitDot:       config.ImplicitDot,
		implicitDotIgnore: utils.NewSetFromSlice(config.ImplicitDotIgnore),
	}
	for _, sourceFile := range sourceFiles {
		mapping.Add(sourceFile)
	}
	return mapping
}

func (fm *FileMapping) Add(newSource RelativePath) {
	target := fm.targetBaseDir.JoinPath(fm.mapSourceToTarget(newSource))
	if existingSource, contains := fm.mapping[target]; contains {
		log.Warning("Conflicting files: %s and %s both map to %s. Ignoring %s", existingSource, newSource, target, newSource)
	} else {
		fm.mapping[target] = fm.sourceBaseDir.JoinPath(newSource)
	}
}

func (fm *FileMapping) GetTargets() []AbsolutePath {
	targets := make([]AbsolutePath, 0, len(fm.mapping))
	for target := range fm.mapping {
		targets = append(targets, target)
	}
	return targets
}

func (fm *FileMapping) InstallNewLinks(ignore []AbsolutePath) {
	for target, source := range fm.mapping {
		if contains(ignore, target) {
			log.Info("Target %s already exists and will not be created", target)
		} else {
			log.Info("Linking %s -> %s", target, source)
			err := os.Symlink(source.Str(), target.Str())
			if err != nil {
				handleSymlinkError(target, source, err)
			}
		}
	}
}

func handleSymlinkError(target, source AbsolutePath, err error) {
	// The most likely reason os.Symlink failed is that the target (symlink path) already exists
	fileInfo, statErr := os.Lstat(target.Str())
	if statErr != nil {
		// Either the target does not exist or we cannot access it, either way this is unexpected
		log.Error("Failed to create link %s -> %s: %s", target, source, err)
		return
	}

	isSymlink := fileInfo.Mode()&os.ModeSymlink != 0
	if isSymlink {
		// Check if the existing symlink points to the correct source
		linkSource, linkErr := os.Readlink(target.Str())
		if linkErr != nil {
			log.Error("Failed to read link %s: %s", target, linkErr)
			return
		}
		if linkSource == source.Str() {
			log.Info("Link %s -> %s already existed, cache was incorrect!", target, source)
			return
		}
	}

	isRegularFile := fileInfo.Mode().IsRegular()
	if isRegularFile {
		log.Warning("Failed to create link %s -> %s: target already exists", target, source)
	} else {
		log.Error("Failed to create link %s -> %s: target is not a symlink", target, source)
	}
}

func (fm *FileMapping) RemoveStaleLinks(previousTargets []AbsolutePath) {
	for _, previousTarget := range previousTargets {
		if _, contains := fm.mapping[previousTarget]; !contains {
			log.Info("Removing stale link %s", previousTarget)
		}
	}
}

func contains[T comparable](slice []T, element T) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}

func (fm *FileMapping) mapSourceToTarget(source RelativePath) RelativePath {
	target := source
	if fm.implicitDot && !fm.implicitDotIgnore.Contains(getTopLevelDir(source)) {
		target = "." + source
	}
	return target
}

func getTopLevelDir(filePath RelativePath) string {
	filePathStr := string(filePath)
	firstSeparatorIndex := strings.IndexRune(filePathStr, filepath.Separator)
	if firstSeparatorIndex == -1 {
		return filePathStr
	}
	return filePathStr[:firstSeparatorIndex]
}
