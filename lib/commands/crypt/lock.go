package crypt

import (
	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/log"
	"github.com/pol-rivero/doot/lib/utils"
)

func Lock(force bool) {
	ensureGitCryptInstalled()
	dotfilesDir := common.FindDotfilesDir()
	ensureGitCryptIsInitialized(dotfilesDir)

	var err error
	if force {
		err = utils.RunCommand(dotfilesDir, "git-crypt", "lock", "--force")
	} else {
		err = utils.RunCommand(dotfilesDir, "git-crypt", "lock")
	}
	if err != nil {
		log.Fatal("Failed to lock repository, make sure your working directory is clean or use 'doot crypt lock --force'.")
	}

	log.Printlnf("Repository locked successfully.")
}
