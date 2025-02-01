package helpers

import "log"

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
		log.Printf("Conflicting files: %s and %s both map to %s. Ignoring %s", existingSource, newSource, target, newSource)
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
			log.Printf("Target %s already exists and will not be created", target)
		} else {
			log.Printf("Linking %s -> %s", target, source)
		}
	}
}

func (fm *FileMapping) RemoveStaleLinks(previousTargets []string) {
	for _, previousTarget := range previousTargets {
		if _, contains := fm.mapping[previousTarget]; !contains {
			log.Printf("Removing stale link %s", previousTarget)
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
