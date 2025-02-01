package helpers

import (
	"github.com/pol-rivero/doot/lib/log"
)

type FileMapping struct {
	// Map target -> source
	mapping map[string]string
}

func NewFileMapping(capacity int) FileMapping {
	return FileMapping{
		mapping: make(map[string]string, capacity),
	}
}

func (fm *FileMapping) Add(newSource string, target string) {
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
