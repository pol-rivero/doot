package install

import (
	"os"
	"os/exec"
	"slices"
	"strings"

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
	mapping           map[AbsolutePath]SourcePath // Map target (symlink) -> source (dotfile)
	sourceBaseDir     AbsolutePath
	targetBaseDir     AbsolutePath
	implicitDot       bool
	implicitDotIgnore set.Set[string]
	targetsSkipped    []AbsolutePath
}

func NewFileMapping(dotfilesDir AbsolutePath, config *config.Config, sourceFiles []RelativePath) FileMapping {
	mapping := FileMapping{
		mapping:           make(map[AbsolutePath]SourcePath, len(sourceFiles)),
		sourceBaseDir:     dotfilesDir,
		targetBaseDir:     NewAbsolutePath(config.TargetDir),
		implicitDot:       config.ImplicitDot,
		implicitDotIgnore: set.NewFromSlice(config.ImplicitDotIgnore),
		targetsSkipped:    make([]AbsolutePath, 0),
	}
	for _, sourceFile := range sourceFiles {
		mapping.Add(sourceFile)
	}
	return mapping
}

func (fm *FileMapping) Add(newSource RelativePath) {
	relativeTarget, newIsHostSpecific := fm.mapSourceToTarget(newSource)
	if !relativeTarget.HasValue() {
		return
	}
	target := fm.targetBaseDir.JoinPath(relativeTarget.Value())

	oldSource, oldSourceExists := fm.mapping[target]
	preferNewSource := !oldSourceExists || newIsHostSpecific
	if preferNewSource {
		fm.mapping[target] = SourcePath{
			path:         fm.sourceBaseDir.JoinPath(newSource),
			hostSpecific: newIsHostSpecific,
		}
		if oldSourceExists {
			log.Info("Host-specific file %s overrides %s for target %s", newSource, oldSource.path, target)
		}
	} else if oldSource.hostSpecific {
		log.Info("Host-specific file %s overrides %s for target %s", oldSource.path, newSource, target)
	} else if newIsHostSpecific {
		// This is rare, but it can happen if 2 files map to the same target after removing '.doot-crypt' or adding the implicit dot
		log.Warning("Conflicting files: %s and %s both map to %s. Ignoring %s", oldSource.path, newSource, target, newSource)
	}
}

func (fm *FileMapping) GetInstalledTargets() []AbsolutePath {
	targets := make([]AbsolutePath, 0, len(fm.mapping))
	for target := range fm.mapping {
		if !slices.Contains(fm.targetsSkipped, target) {
			targets = append(targets, target)
		}
	}
	return targets
}

func (fm *FileMapping) InstallNewLinks(ignore []AbsolutePath) int {
	createdLinksCount := 0
	for target, sourceStruct := range fm.mapping {
		source := sourceStruct.path
		if slices.Contains(ignore, target) {
			log.Info("Target %s already exists and will not be created", target)
			continue
		}
		fileInfo, err := os.Lstat(target.Str())
		if err == nil {
			added := fm.handleTargetAlreadyExists(fileInfo, target, source)
			if added {
				createdLinksCount++
			}
			continue
		}
		if os.IsNotExist(err) && EnsureParentDir(target) {
			log.Info("Linking %s -> %s", target, source)
			err = os.Symlink(source.Str(), target.Str())
			if err == nil {
				createdLinksCount++
				continue
			}
		}
		log.Error("Failed to create link %s -> %s: %s", target, source, err)
	}
	return createdLinksCount
}

func (fm *FileMapping) RemoveStaleLinks(previousTargets []AbsolutePath) int {
	removedLinksCount := 0
	for _, previousTarget := range previousTargets {
		if _, contains := fm.mapping[previousTarget]; !contains {
			if !canBeSafelyRemoved(previousTarget, fm.sourceBaseDir) {
				log.Info("%s appears to have been modified externally. Skipping removal to avoid data loss.", previousTarget)
				continue
			}
			log.Info("Removing link %s", previousTarget)
			success := RemoveAndCleanup(previousTarget, fm.targetBaseDir)
			if success {
				removedLinksCount++
			}
		}
	}
	return removedLinksCount
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
			printDiff(target, source)
		}
	}
}

func (fm *FileMapping) mapSourceToTarget(source RelativePath) (optional.Optional[RelativePath], bool) {
	target := source
	// The doot directory should not be symlinked
	if strings.HasPrefix(source.Str(), "doot/") {
		return optional.Empty[RelativePath](), false
	}
	if fm.implicitDot && !fm.implicitDotIgnore.Contains(source.TopLevelDir()) && !strings.HasPrefix(source.Str(), ".") {
		target = "." + source
	}
	target = target.Replace(".doot-crypt", "")
	return optional.Of(target), false
}

func canBeSafelyRemoved(linkPath AbsolutePath, expectedDestinationDir AbsolutePath) bool {
	linkSource, linkErr := os.Readlink(linkPath.Str())
	if linkErr != nil {
		return false
	}
	return strings.HasPrefix(linkSource, expectedDestinationDir.Str())
}

func printDiff(leftFile AbsolutePath, rightFile AbsolutePath) {
	cmd := exec.Command("diff", "-u", leftFile.Str(), rightFile.Str())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Error("Failed to run diff: %s", err)
	}
}
