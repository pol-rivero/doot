package install

import (
	"os"
	"slices"
	"strings"

	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/config"
	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/pol-rivero/doot/lib/utils"
	"github.com/pol-rivero/doot/lib/utils/optional"
	"github.com/pol-rivero/doot/lib/utils/set"
)

type SourcePath struct {
	path         AbsolutePath
	hostSpecific bool
}

type FileMapping struct {
	mapping           map[AbsolutePath]SourcePath // Installed target (symlink path) -> source (dotfile, symlink content/target)
	sourceBaseDir     AbsolutePath
	targetBaseDir     AbsolutePath
	implicitDot       bool
	implicitDotIgnore set.Set[string]
	hostnameFilter    HostnameFilter
	diffCommand       string
	targetsSkipped    []AbsolutePath
}

func NewFileMapping(dotfilesDir AbsolutePath, config *config.Config, sourceFiles []RelativePath) FileMapping {
	mapping := FileMapping{
		mapping:           make(map[AbsolutePath]SourcePath, len(sourceFiles)),
		sourceBaseDir:     dotfilesDir,
		targetBaseDir:     NewAbsolutePath(config.TargetDir),
		implicitDot:       config.ImplicitDot,
		implicitDotIgnore: set.NewFromSlice(config.ImplicitDotIgnore),
		hostnameFilter:    getHostnameFilter(config.Hosts),
		diffCommand:       config.DiffCommand,
		targetsSkipped:    make([]AbsolutePath, 0),
	}
	for _, sourceFile := range sourceFiles {
		mapping.Add(sourceFile)
	}
	return mapping
}

func (fm *FileMapping) Add(relativeSource RelativePath) {
	relativeTarget, newIsHostSpecific := fm.mapSourceToTarget(relativeSource)
	if relativeTarget.IsEmpty() {
		return
	}
	source := fm.sourceBaseDir.JoinPath(relativeSource)
	target := fm.targetBaseDir.JoinPath(relativeTarget.Value())

	oldSource, oldSourceExists := fm.mapping[target]
	preferNewSource := !oldSourceExists || newIsHostSpecific
	if preferNewSource {
		fm.mapping[target] = SourcePath{
			path:         source,
			hostSpecific: newIsHostSpecific,
		}
		if oldSourceExists {
			log.Info("Host-specific file %s overrides %s for target %s", source, oldSource.path, target)
		}
	} else if oldSource.hostSpecific {
		log.Info("Host-specific file %s overrides %s for target %s", oldSource.path, source, target)
	} else {
		// This is rare, but it can happen if 2 files map to the same target after removing '.doot-crypt' or adding the implicit dot
		log.Warning("Conflicting files: %s and %s both map to %s. Ignoring %s", oldSource.path, source, target, source)
	}
}

func (fm *FileMapping) GetInstalledTargets() SymlinkCollection {
	targets := NewSymlinkCollection(len(fm.mapping))
	for targetPath, sourcePath := range fm.mapping {
		if !slices.Contains(fm.targetsSkipped, targetPath) {
			targets.Add(targetPath, sourcePath.path)
		}
	}
	return targets
}

func (fm *FileMapping) InstallNewLinks(alreadyExist *SymlinkCollection) []AbsolutePath {
	createdLinks := make([]AbsolutePath, 0, 5)
	for target, sourceStruct := range fm.mapping {
		newSource := sourceStruct.path
		oldSource := alreadyExist.Get(target)
		if oldSource.HasValue() {
			added := fm.handleOutdatedLink(target, oldSource.Value(), newSource)
			if added {
				createdLinks = append(createdLinks, target)
			}
			continue
		}
		fileInfo, err := os.Lstat(target.Str())
		if err == nil {
			added := fm.handleTargetAlreadyExists(fileInfo, target, newSource)
			if added {
				createdLinks = append(createdLinks, target)
			}
			continue
		}
		if os.IsNotExist(err) && EnsureParentDir(target) {
			log.Info("Linking %s -> %s", target, newSource)
			err = os.Symlink(newSource.Str(), target.Str())
			if err == nil {
				createdLinks = append(createdLinks, target)
				continue
			}
		}
		log.Error("Failed to create link %s -> %s: %s", target, newSource, err)
	}
	return createdLinks
}

func (fm *FileMapping) RemoveStaleLinks(previousLinks *SymlinkCollection) []AbsolutePath {
	removedLinks := make([]AbsolutePath, 0, 5)
	for previousLinkPath := range previousLinks.Iter() {
		if _, contains := fm.mapping[previousLinkPath]; !contains {
			if !canBeSafelyRemoved(previousLinkPath, fm.sourceBaseDir) {
				log.Info("%s appears to have been modified externally. Skipping removal to avoid data loss.", previousLinkPath)
				continue
			}
			log.Info("Removing link %s", previousLinkPath)
			success := RemoveAndCleanup(previousLinkPath, fm.targetBaseDir)
			if success {
				removedLinks = append(removedLinks, previousLinkPath)
			}
		}
	}
	return removedLinks
}

func (fm *FileMapping) handleTargetAlreadyExists(fileInfo os.FileInfo, target, source AbsolutePath) bool {
	if fileInfo.Mode()&os.ModeSymlink != 0 {
		return fm.handleExistingSymlink(target, source)
	} else if fileInfo.Mode().IsRegular() {
		return fm.handleExistingFile(target, source)
	} else {
		log.Warning("Target %s exists but is not a symlink or a regular file, skipping", target)
		return false
	}
}

func (fm *FileMapping) handleExistingSymlink(target, source AbsolutePath) bool {
	linkSource, linkErr := os.Readlink(target.Str())
	if linkErr != nil {
		log.Error("Failed to read link %s: %s", target, linkErr)
		return false
	}
	if linkSource == source.Str() {
		log.Info("Link %s -> %s already existed (cache was outdated)", target, source)
		return false
	}
	replace := utils.RequestInput("yN", "Link %s already exists, but it points to %s instead of %s. Replace it?", target, linkSource, source)
	if replace == 'y' {
		err := ReplaceWithSymlink(target, source)
		return err == nil
	} else {
		fm.targetsSkipped = append(fm.targetsSkipped, target)
		return false
	}
}

func (fm *FileMapping) handleOutdatedLink(target, oldSource, newSource AbsolutePath) bool {
	if oldSource == newSource {
		log.Info("Target %s already exists and will not be created", target)
		return false
	} else {
		log.Info("Link %s is outdated (was %s, changing to %s)", target, oldSource, newSource)
		err := ReplaceWithSymlink(target, newSource)
		return err == nil
	}
}

func (fm *FileMapping) handleExistingFile(target, source AbsolutePath) bool {
	contents, readErr := os.ReadFile(target.Str())
	if readErr != nil {
		log.Error("Failed to read file %s: %s", target, readErr)
		return false
	}
	sourceContents, readErr := os.ReadFile(source.Str())
	if readErr != nil {
		log.Error("Failed to read file %s: %s", source, readErr)
		return false
	}
	if string(contents) == string(sourceContents) {
		log.Info("File %s exists but its contents are identical to %s, replacing silently", target, source)
		err := ReplaceWithSymlink(target, source)
		return err == nil
	}
	for {
		replace := utils.RequestInput("yNd", "File %s already exists, but its contents differ from %s. Replace it? (press D to see diff)", target, source)
		switch replace {
		case 'y':
			err := ReplaceWithSymlink(target, source)
			return err == nil
		case 'n':
			fm.targetsSkipped = append(fm.targetsSkipped, target)
			return false
		case 'd':
			fm.printDiff(target, source)
		}
	}
}

func (fm *FileMapping) mapSourceToTarget(source RelativePath) (optional.Optional[RelativePath], bool) {
	target := source
	if fm.hostnameFilter.isIgnored(source) {
		return optional.Empty[RelativePath](), false
	}
	isHostSpecific, prefixLen := fm.hostnameFilter.isHostSpecific(source)
	if isHostSpecific {
		target = target.RemoveBaseDir(prefixLen)
	}
	if fm.implicitDot && !fm.implicitDotIgnore.Contains(source.TopLevelDir()) && !strings.HasPrefix(target.Str(), ".") {
		target = "." + target
	}
	target = target.Replace(common.DOOT_CRYPT_EXT, "")
	return optional.WrapString(target), isHostSpecific
}

func canBeSafelyRemoved(linkPath AbsolutePath, expectedDestinationDir AbsolutePath) bool {
	linkSource, linkErr := os.Readlink(linkPath.Str())
	if linkErr != nil {
		return false
	}
	return strings.HasPrefix(linkSource, expectedDestinationDir.Str())
}

func (fm *FileMapping) printDiff(leftFile AbsolutePath, rightFile AbsolutePath) {
	err := utils.RunCommandStr(fm.sourceBaseDir, fm.diffCommand, leftFile.Str(), rightFile.Str())
	if err != nil {
		log.Info("Diff command had non-zero exit code: %s. This is usually not a problem.", err)
	}
}
