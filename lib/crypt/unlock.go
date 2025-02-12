package crypt

import (
	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/pol-rivero/doot/lib/utils"
	"github.com/pol-rivero/doot/lib/utils/optional"
)

func Unlock(keyFile optional.Optional[string]) {
	ensureGitCryptInstalled()
	dotfilesDir := common.FindDotfilesDir()

	if gitAttributesIsSet(dotfilesDir) {
		log.Info("Git attributes already set, skipping...")
	} else {
		err := appendGitAttributes(dotfilesDir)
		if err != nil {
			log.Fatal("Failed to append git attributes: %s", err)
		}
	}

	if keyFile.HasValue() {
		unlockWithKeyFile(dotfilesDir, keyFile.Value())
	} else {
		unlockGPG(dotfilesDir)
	}

	log.Printlnf("Repository (%s) unlocked successfully, you now have access to the encrypted files", dotfilesDir)
}

func unlockWithKeyFile(dotfilesDir AbsolutePath, keyFile string) {
	err := utils.RunCommand(dotfilesDir, "git-crypt", "unlock", keyFile)
	if err != nil {
		log.Fatal("Failed to unlock repository")
	}
}

func unlockGPG(dotfilesDir AbsolutePath) {
	err := utils.RunCommand(dotfilesDir, "git-crypt", "unlock")
	if err != nil {
		log.Fatal("Failed to unlock repository with GPG: %s", err)
	}
}
