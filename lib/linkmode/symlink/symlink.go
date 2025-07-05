package linkmode_symlink

import (
	"os"
	"strings"

	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
)

type SymlinkLinkMode struct{}

func (l *SymlinkLinkMode) CreateLink(dotfilesSource, target AbsolutePath) error {
	return os.Symlink(dotfilesSource.Str(), target.Str())
}

func (l *SymlinkLinkMode) IsInstalledLinkOf(maybeInstalledLinkPath string, dotfilePath AbsolutePath) bool {
	fileInfo, err := os.Lstat(maybeInstalledLinkPath)
	if err != nil {
		log.Info("Failed to stat %s: %v", maybeInstalledLinkPath, err)
		return false
	}
	if fileInfo.Mode()&os.ModeSymlink != 0 && getSymlinkTarget(maybeInstalledLinkPath) == dotfilePath.Str() {
		return true
	}
	return false
}

func getSymlinkTarget(linkPath string) string {
	linkSource, linkErr := os.Readlink(linkPath)
	if linkErr != nil {
		log.Fatal("Failed to read link %s: %v", linkPath, linkErr)
	}
	return linkSource
}

func (l *SymlinkLinkMode) CanBeSafelyRemoved(linkPath AbsolutePath, expectedDestinationDir string) bool {
	linkSource, linkErr := os.Readlink(linkPath.Str())
	if linkErr != nil {
		return false
	}
	return strings.HasPrefix(linkSource, expectedDestinationDir)
}
