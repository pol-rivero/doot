package linkmode

import (
	"fmt"
	"os"
	"syscall"

	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
)

type HardlinkLinkMode struct{}

func (l *HardlinkLinkMode) CreateLink(dotfilesSource, target AbsolutePath) error {
	return os.Link(dotfilesSource.Str(), target.Str())
}

func (l *HardlinkLinkMode) IsInstalledLinkOf(maybeInstalledLinkPath string, dotfilePath AbsolutePath) bool {
	info1, err1 := unixStat(maybeInstalledLinkPath)
	if err1 != nil {
		log.Info("Failed to stat %s: %v", maybeInstalledLinkPath, err1)
		return false
	}
	info2, err2 := unixStat(string(dotfilePath))
	if err2 != nil {
		log.Info("Failed to stat %s: %v", dotfilePath, err2)
		return false
	}
	return info1.Dev == info2.Dev && info1.Ino == info2.Ino
}

func (l *HardlinkLinkMode) CanBeSafelyRemoved(linkPath AbsolutePath, _ string) bool {
	// For hardlinks, we can check for point equality (IsInstalledLinkOf) but we cannot check if the other hardlink is in a given directory.
	// The best we can do is to check if there is another hardlink to the same inode. Since data won't be lost, we can safely remove it.
	info, err := unixStat(linkPath.Str())
	if err != nil {
		log.Error("Failed to stat %s: %v", linkPath, err)
		return false
	}
	if info.Nlink > 1 {
		log.Info("Link %s has %d hardlinks, can be safely removed", linkPath, info.Nlink)
		return true
	}
	return false
}

func unixStat(path string) (*syscall.Stat_t, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, fmt.Errorf("failed to get Stat_t for %s (hardlink_mode is not supported on Windows)", path)
	}
	return stat, nil
}
