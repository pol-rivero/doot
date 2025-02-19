package crypt

import (
	"path/filepath"

	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/log"
	"github.com/pol-rivero/doot/lib/utils"
)

func ExportKey(outputPath string) {
	ensureGitCryptInstalled()
	dotfilesDir := common.FindDotfilesDir()
	ensureGitCryptIsInitialized(dotfilesDir)

	outputPath, err := filepath.Abs(outputPath)
	if err != nil {
		log.Fatal("Failed to get absolute path for '%s': %v", outputPath, err)
	}

	err = utils.RunCommand(dotfilesDir, "git-crypt", "export-key", outputPath)
	if err != nil {
		log.Fatal("Failed to export key")
	}

	log.Printlnf("Key exported to '%s'", outputPath)
}
