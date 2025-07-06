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
	// For hardlinks, we can check for point equality (IsInstalledLinkOf) but we cannot check if the other hardlink is in a given directory.
	// The best we can do is to check if there is another hardlink to the same inode. Since data won't be lost, we can safely remove it.
	info, err := osStat(linkPath.Str())
	if err != nil {
		log.Error("Failed to stat %s: %v", linkPath, err)
		return false
	}
	if info.numLinks > 1 {
		log.Info("Link %s has %d hardlinks, can be safely removed", linkPath, info.numLinks)
		return true
	}
	return false
}
