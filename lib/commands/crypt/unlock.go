package crypt

import (
	"errors"
	"path/filepath"

	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/pol-rivero/doot/lib/utils"
	"github.com/pol-rivero/doot/lib/utils/optional"
)

func Unlock(keyFile optional.Optional[string]) {
	dotfilesDir := common.FindDotfilesDir()
	err := UnlockOrErr(dotfilesDir, keyFile)
	if err != nil {
		log.Error("Unlock error: %s", err)
	} else {
		log.Printlnf("Repository (%s) unlocked successfully, you now have access to the encrypted files", dotfilesDir)
	}
}

func UnlockOrErr(dotfilesDir AbsolutePath, keyFile optional.Optional[string]) error {
	ensureGitCryptInstalled()

	if gitAttributesIsSet(dotfilesDir) {
		log.Info("Git attributes already set, skipping...")
	} else {
		err := appendGitAttributes(dotfilesDir)
		if err != nil {
			log.Fatal("Failed to append git attributes: %s", err)
		}
	}

	if keyFile.HasValue() {
		return unlockWithKeyFile(dotfilesDir, keyFile.Value())
	} else {
		return unlockGPG(dotfilesDir)
	}
}

func unlockWithKeyFile(dotfilesDir AbsolutePath, keyFile string) error {
	keyFile, err := filepath.Abs(keyFile)
	if err != nil {
		log.Fatal("Failed to get absolute path for '%s': %v", keyFile, err)
	}

	err = utils.RunCommand(dotfilesDir, "git-crypt", "unlock", keyFile)
	if err != nil {
		return errors.New("failed to unlock repository, make sure the provided key file is correct")
	}
	return nil
}

func unlockGPG(dotfilesDir AbsolutePath) error {
	err := utils.RunCommand(dotfilesDir, "git-crypt", "unlock")
	if err != nil {
		return errors.New("failed to unlock repository with GPG, make sure you have the correct private key installed")
	}
	return nil
}
