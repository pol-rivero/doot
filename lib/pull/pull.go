package pull

import (
	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/log"
	"github.com/pol-rivero/doot/lib/install"
	"github.com/pol-rivero/doot/lib/utils"
)

func Pull() {
	dotfilesDir := common.FindDotfilesDir()

	err := utils.RunCommand(dotfilesDir, "git", "pull", "--recurse-submodules")
	if err != nil {
		log.Fatal("git pull failed, check the error above for more information.")
	}

	log.Info("Changes pulled successfully. Proceeding to install.")
	install.Install(false)
}
