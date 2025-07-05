package linkmode

import (
	"github.com/pol-rivero/doot/lib/common/config"
	. "github.com/pol-rivero/doot/lib/types"
)

type LinkMode interface {
	CreateLink(dotfilesSource, target AbsolutePath) error
	IsInstalledLinkOf(maybeInstalledLinkPath string, dotfilePath AbsolutePath) bool
	CanBeSafelyRemoved(linkPath AbsolutePath, expectedDestinationDir string) bool
}

func GetLinkMode(config *config.Config) LinkMode {
	if config.HardlinkMode {
		return &HardlinkLinkMode{}
	}
	return &SymlinkLinkMode{}
}
