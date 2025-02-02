package helpers

import (
	"path/filepath"

	"github.com/pol-rivero/doot/lib/log"
	"github.com/pol-rivero/glob"
)

type GlobCollection struct {
	globs []glob.Glob
}

func NewGlobCollection(patterns []string) GlobCollection {
	globs := make([]glob.Glob, 0, len(patterns))
	for _, pattern := range patterns {
		g, err := glob.Compile(pattern, filepath.Separator)
		if err != nil {
			log.Warning("Ignoring invalid glob pattern '%s': %v", pattern, err)
			continue
		}
		globs = append(globs, g)
	}
	return GlobCollection{globs}
}

func (gc *GlobCollection) Matches(s string) bool {
	for _, g := range gc.globs {
		if g.Match(s) {
			return true
		}
	}
	return false
}

func (gc *GlobCollection) Len() int {
	return len(gc.globs)
}
