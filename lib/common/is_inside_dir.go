package common

import (
	"path/filepath"
	"strings"
)

// Reports whether path is strictly inside dir.
func IsInsideDir(dir, path string) bool {
	rel, err := filepath.Rel(dir, path)
	if err != nil || filepath.IsAbs(rel) {
		return false
	}
	return rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
