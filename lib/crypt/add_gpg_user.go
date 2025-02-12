package crypt

import (
	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/log"
	"github.com/pol-rivero/doot/lib/utils"
)

func AddGpgUser(userId string) {
	ensureGitCryptInstalled()
	dotfilesDir := common.FindDotfilesDir()
	ensureGitCryptIsInitialized(dotfilesDir)

	err := utils.RunCommand(dotfilesDir, "git-crypt", "add-gpg-user", userId)
	if err != nil {
		log.Fatal("Failed to add GPG user")
	}

	log.Printlnf("GPG user '%s' added successfully.", userId)
}
