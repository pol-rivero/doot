package install

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/pol-rivero/doot/lib/config"
	"github.com/pol-rivero/doot/lib/log"
	"github.com/pol-rivero/doot/lib/utils"
)

type FileMapping struct {
	// Map target -> source
	mapping           map[string]string
	targetBaseDir     string
	implicitDot       bool
	implicitDotIgnore utils.Set[string]
}

func NewFileMapping(config *config.Config, sourceFiles []string) FileMapping {
	mapping := FileMapping{
		mapping:           make(map[string]string, len(sourceFiles)),
		targetBaseDir:     config.TargetDir,
		implicitDot:       config.ImplicitDot,
		implicitDotIgnore: utils.NewSetFromSlice(config.ImplicitDotIgnore),
	}
	for _, sourceFile := range sourceFiles {
		mapping.Add(sourceFile)
	}
	return mapping
}

func (fm *FileMapping) Add(newSource string) {
	target := path.Join(fm.targetBaseDir, fm.mapSourceToTarget(newSource))
	if existingSource, contains := fm.mapping[target]; contains {
		log.Warning("Conflicting files: %s and %s both map to %s. Ignoring %s", existingSource, newSource, target, newSource)
	} else {
		fm.mapping[target] = newSource
	}
}

func (fm *FileMapping) GetTargets() []string {
	targets := make([]string, 0, len(fm.mapping))
	for target := range fm.mapping {
		targets = append(targets, target)
	}
	return targets
}

func (fm *FileMapping) InstallNewLinks(ignore []string) {
	for target, source := range fm.mapping {
		if contains(ignore, target) {
			log.Info("Target %s already exists and will not be created", target)
		} else {
			log.Info("Linking %s -> %s", target, source)
		}
	}
}

func (fm *FileMapping) RemoveStaleLinks(previousTargets []string) {
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

func (fm *FileMapping) mapSourceToTarget(sourceRelativePath string) string {
	target := sourceRelativePath
	if fm.implicitDot && !fm.implicitDotIgnore.Contains(getTopLevelDir(sourceRelativePath)) {
		target = "." + sourceRelativePath
	}
	return target
}

func getTopLevelDir(filePath string) string {
	firstSeparatorIndex := strings.IndexRune(filePath, filepath.Separator)
	if firstSeparatorIndex == -1 {
		return filePath
	}
	return filePath[:firstSeparatorIndex]
}
