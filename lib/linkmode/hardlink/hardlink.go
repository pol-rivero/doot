package linkmode_hardlink

import (
	"os"

	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
)

type HardlinkLinkMode struct{}

type OsStatResult struct {
	numLinks   uint64
	hardlinkId HardlinkId
}

func (l *HardlinkLinkMode) CreateLink(dotfilesSource, target AbsolutePath) error {
	return os.Link(dotfilesSource.Str(), target.Str())
}

func (l *HardlinkLinkMode) IsInstalledLinkOf(maybeInstalledLinkPath string, dotfilePath AbsolutePath) bool {
	info1, err1 := osStat(maybeInstalledLinkPath)
	if err1 != nil {
		log.Info("Failed to stat %s: %v", maybeInstalledLinkPath, err1)
		return false
	}
	info2, err2 := osStat(string(dotfilePath))
	if err2 != nil {
		log.Info("Failed to stat %s: %v", dotfilePath, err2)
		return false
	}
	return info1.hardlinkId == info2.hardlinkId
}

func (l *HardlinkLinkMode) CanBeSafelyRemoved(linkPath AbsolutePath, _ string) bool {
	// Hardlinks are just names for the an inode, so we cannot check if this is the same inode we installed without storing
	// the HardlinkId in the cache. Doing that would greatly increase complexity and it's probably not worth it, since the
	// probability of actual data loss is low and I doubt many people will use hardlinks.
	return true
}
