package bootstrap

import (
	"os"

	"github.com/pol-rivero/doot/lib/commands/install"
	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/pol-rivero/doot/lib/utils/optional"
)

func Bootstrap(repo string, dotfilesDirRel string, keyPath optional.Optional[string]) {
	dotfilesDir := RelativeToPWD(dotfilesDirRel)
	CloneRepoOrExit(repo, dotfilesDir)

	err := os.Setenv(common.ENV_DOOT_DIR, dotfilesDir.Str())
	if err != nil {
		log.Fatal("Error setting DOOT_DIR environment variable: %v", err)
	}

	UnlockIfNeeded(dotfilesDir, keyPath)

	common.RunHooks(dotfilesDir, "before-bootstrap")
	install.Install(false)
	common.RunHooks(dotfilesDir, "after-bootstrap")
}
