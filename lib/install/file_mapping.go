package install

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pol-rivero/doot/lib/config"
	"github.com/pol-rivero/doot/lib/log"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/pol-rivero/doot/lib/utils"
	"github.com/pol-rivero/doot/lib/utils/optional"
)

type FileMapping struct {
	mapping           map[AbsolutePath]AbsolutePath // Map target (symlink) -> source (dotfile)
	sourceBaseDir     AbsolutePath
	targetBaseDir     AbsolutePath
	implicitDot       bool
	implicitDotIgnore utils.Set[string]
	targetsSkipped    []AbsolutePath
}

func NewFileMapping(dotfilesDir AbsolutePath, config *config.Config, sourceFiles []RelativePath) FileMapping {
	mapping := FileMapping{
		mapping:           make(map[AbsolutePath]AbsolutePath, len(sourceFiles)),
		sourceBaseDir:     dotfilesDir,
		targetBaseDir:     NewAbsolutePath(config.TargetDir),
		implicitDot:       config.ImplicitDot,
		implicitDotIgnore: utils.NewSetFromSlice(config.ImplicitDotIgnore),
		targetsSkipped:    make([]AbsolutePath, 0),
	}
	for _, sourceFile := range sourceFiles {
		mapping.Add(sourceFile)
	}
	return mapping
}

func (fm *FileMapping) Add(newSource RelativePath) {
	relativeTarget := fm.mapSourceToTarget(newSource)
	if !relativeTarget.HasValue() {
		return
	}
	target := fm.targetBaseDir.JoinPath(relativeTarget.Value())
	if existingSource, contains := fm.mapping[target]; contains {
		log.Warning("Conflicting files: %s and %s both map to %s. Ignoring %s", existingSource, newSource, target, newSource)
	} else {
		fm.mapping[target] = fm.sourceBaseDir.JoinPath(newSource)
	}
}

func (fm *FileMapping) GetInstalledTargets() []AbsolutePath {
	targets := make([]AbsolutePath, 0, len(fm.mapping))
	for target := range fm.mapping {
		if !contains(fm.targetsSkipped, target) {
			targets = append(targets, target)
		}
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
				fm.handleSymlinkError(target, source)
			}
		}
	}
}

func (fm *FileMapping) RemoveStaleLinks(previousTargets []AbsolutePath) {
	for _, previousTarget := range previousTargets {
		if _, contains := fm.mapping[previousTarget]; !contains {
			log.Info("Removing stale link %s", previousTarget)
			err := os.Remove(previousTarget.Str())
			if err != nil {
				log.Error("Failed to remove stale link %s: %s", previousTarget, err)
			}
		}
	}
}

func (fm *FileMapping) handleSymlinkError(target, source AbsolutePath) {
	// The most likely reason os.Symlink failed is that the target (symlink path) already exists
	fileInfo, err := os.Lstat(target.Str())
	if err != nil {
		// Either the target does not exist or we cannot access it, either way this is unexpected
		log.Error("Failed to create link %s -> %s: %s", target, source, err)
		return
	}

	if fileInfo.Mode()&os.ModeSymlink != 0 {
		fm.handleExistingSymlink(target, source)
	} else if fileInfo.Mode().IsRegular() {
		fm.handleExistingFile(target, source)
	} else {
		log.Error("Failed to create link %s -> %s: %s", target, source, err)
	}
}

func (fm *FileMapping) handleExistingSymlink(target, source AbsolutePath) {
	linkSource, linkErr := os.Readlink(target.Str())
	if linkErr != nil {
		log.Error("Failed to read link %s: %s", target, linkErr)
		return
	}
	if linkSource == source.Str() {
		log.Info("Link %s -> %s already existed (cache was outdated)", target, source)
		return
	}
	replace := utils.RequestInput("yN", "Link %s already exists, but points to %s instead of %s. Replace it?", target, linkSource, source)
	if replace == 'y' {
		err := utils.ReplaceWithSymlink(target, source)
		if err != nil {
			return
		}
	} else {
		fm.targetsSkipped = append(fm.targetsSkipped, target)
	}
}

func (fm *FileMapping) handleExistingFile(target, source AbsolutePath) {
	contents, readErr := os.ReadFile(target.Str())
	if readErr != nil {
		log.Error("Failed to read file %s: %s", target, readErr)
		return
	}
	sourceContents, readErr := os.ReadFile(source.Str())
	if readErr != nil {
		log.Error("Failed to read file %s: %s", source, readErr)
		return
	}
	if string(contents) == string(sourceContents) {
		log.Info("File %s exists but its contents are identical to %s, replacing silently", target, source)
		err := utils.ReplaceWithSymlink(target, source)
		if err != nil {
			return
		}
	}
	replace := ' '
	for replace != 'y' && replace != 'n' {
		replace = utils.RequestInput("yNd", "File %s already exists, but its contents differ from %s. Replace it? (press D to see diff)", target, source)
		switch replace {
		case 'y':
			err := utils.ReplaceWithSymlink(target, source)
			if err != nil {
				return
			}
		case 'n':
			fm.targetsSkipped = append(fm.targetsSkipped, target)
		case 'd':
			printDiff(target, source)
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

func (fm *FileMapping) mapSourceToTarget(source RelativePath) optional.Optional[RelativePath] {
	target := source
	// The doot directory should not be symlinked
	if strings.HasPrefix(source.Str(), "doot/") {
		return optional.Empty[RelativePath]()
	}
	if fm.implicitDot && !fm.implicitDotIgnore.Contains(getTopLevelDir(source)) && !strings.HasPrefix(source.Str(), ".") {
		target = "." + source
	}
	target = target.Replace(".doot-crypt", "")
	return optional.Of(target)
}

func getTopLevelDir(filePath RelativePath) string {
	filePathStr := string(filePath)
	firstSeparatorIndex := strings.IndexRune(filePathStr, filepath.Separator)
	if firstSeparatorIndex == -1 {
		return filePathStr
	}
	return filePathStr[:firstSeparatorIndex]
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
