package glob_collection

import (
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
)

type GlobCollection struct {
	globs []glob.Glob
}

func NewGlobCollection(patterns []string) GlobCollection {
	globs := make([]glob.Glob, 0, len(patterns))
	for _, pattern := range patterns {
		g, err := glob.Compile(preprocessPattern(pattern), filepath.Separator)
		if err != nil {
			log.Warning("Ignoring invalid glob pattern '%s': %v", pattern, err)
			continue
		}
		globs = append(globs, g)
	}
	return GlobCollection{globs}
}

func (gc *GlobCollection) Matches(s RelativePath) bool {
	for _, g := range gc.globs {
		if g.Match(s.Str()) {
			return true
		}
	}
	return false
}

func (gc *GlobCollection) Len() int {
	return len(gc.globs)
}

func preprocessPattern(pattern string) string {
	// Since ** should also match 0-depth directories, we make all instances of **/ optional
	const SUPER_GLOB = "**" + string(filepath.Separator)
	const SUPER_GLOB_REPLACEMENT = "{,**" + string(filepath.Separator) + "}"
	return strings.ReplaceAll(pattern, SUPER_GLOB, SUPER_GLOB_REPLACEMENT)
}
