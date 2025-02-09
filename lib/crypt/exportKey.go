package crypt

import (
	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/log"
	"github.com/pol-rivero/doot/lib/utils"
)

func ExportKey(outputPath string) {
	ensureGitCryptInstalled()
	dotfilesDir := common.FindDotfilesDir()
	ensureGitCryptIsInitialized(dotfilesDir)

	err := utils.RunCommand(dotfilesDir, "git-crypt", "export-key", outputPath)
	if err != nil {
		log.Fatal("Failed to export key: %s", err)
	}

	log.Printf("Key exported to '%s'", outputPath)
}
