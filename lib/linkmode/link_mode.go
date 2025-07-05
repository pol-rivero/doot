package linkmode

import (
	"github.com/pol-rivero/doot/lib/common/cache"
	"github.com/pol-rivero/doot/lib/common/config"
	hardlink "github.com/pol-rivero/doot/lib/linkmode/hardlink"
	symlink "github.com/pol-rivero/doot/lib/linkmode/symlink"
	. "github.com/pol-rivero/doot/lib/types"
)

type LinkMode interface {
	CreateLink(dotfilesSource, target AbsolutePath) error
	IsInstalledLinkOf(maybeInstalledLinkPath string, dotfilePath AbsolutePath) bool
	CanBeSafelyRemoved(linkPath AbsolutePath, expectedDestinationDir string) bool
	RecalculateCache(dotfilesDir AbsolutePath, scanPath string) []*cache.InstalledFile
}

func GetLinkMode(config *config.Config) LinkMode {
	if config.UseHardlinks {
		return &hardlink.HardlinkLinkMode{}
	}
	return &symlink.SymlinkLinkMode{}
}
