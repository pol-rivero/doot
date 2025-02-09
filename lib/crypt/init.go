package crypt

import (
	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/log"
	"github.com/pol-rivero/doot/lib/utils"
)

func Init() {
	ensureGitCryptInstalled()
	dotfilesDir := common.FindDotfilesDir()

	if GitCryptIsInitialized(dotfilesDir) {
		log.Printlnf("Your repository (%s) is already initialized, nothing to do.", dotfilesDir)
		return
	}

	if gitCryptKeyExists(dotfilesDir) {
		log.Info("git-crypt key already exists, skipping...")
	} else {
		err := utils.RunCommand(dotfilesDir, "git-crypt", "init")
		if err != nil {
			log.Fatal("Failed to initialize git-crypt: %s", err)
		}
	}

	if gitAttributesIsSet(dotfilesDir) {
		log.Info("Git attributes already set, skipping...")
	} else {
		err := appendGitAttributes(dotfilesDir)
		if err != nil {
			log.Fatal("Failed to append git attributes: %s", err)
		}
	}

	log.Printlnf("Repository (%s) initialized successfully. Files and directories with '.doot-crypt' will now be encrypted when uploaded to the remote", dotfilesDir)

}
