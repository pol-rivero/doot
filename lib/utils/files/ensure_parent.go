package files

import (
	"os"
	"path/filepath"

	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
)

func EnsureParentDir(target AbsolutePath) bool {
	parentDir := filepath.Dir(target.Str())
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		log.Error("Failed to create directory %s: %s", parentDir, err)
		return false
	}
	return true
}
