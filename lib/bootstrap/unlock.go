package bootstrap

import (
	"os"
	"strings"

	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/log"
	"github.com/pol-rivero/doot/lib/crypt"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/pol-rivero/doot/lib/utils/optional"
)

func UnlockIfNeeded(dotfilesDir AbsolutePath, keyPath optional.Optional[string]) {
	if !containsCryptFiles(dotfilesDir) {
		log.Info("No encrypted files found. Skipping unlock.")
		return
	}
	if !crypt.GitCryptIsInstalled() {
		log.Warning("Git-crypt is not installed. Skipping unlock.")
		return
	}
	if keyPath.IsEmpty() {
		log.Warning("No key provided. Attempting to unlock with GPG...")
	}
	err := crypt.UnlockOrErr(dotfilesDir, keyPath)
	if err != nil {
		log.Warning("Unlock error: %s. Private files won't be installed during bootstrap process.", err)
	}
}

func containsCryptFiles(scanPath AbsolutePath) bool {
	entries, err := os.ReadDir(scanPath.Str())
	if err != nil {
		log.Error("Error reading directory %s: %v", scanPath, err)
		return false
	}
	for _, entry := range entries {
		entryName := entry.Name()
		if strings.Contains(entryName, common.DOOT_CRYPT_EXT) {
			return true
		}
		if entry.IsDir() && containsCryptFiles(scanPath.Join(entryName)) {
			return true
		}
	}
	return false
}
