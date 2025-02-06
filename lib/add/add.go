package add

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/config"
	"github.com/pol-rivero/doot/lib/common/glob_collection"
	"github.com/pol-rivero/doot/lib/common/log"
	"github.com/pol-rivero/doot/lib/install"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/pol-rivero/doot/lib/utils/set"
)

func Add(files []string, crypt bool, hostSpecific bool) {
	dotfilesDir := common.FindDotfilesDir()
	config := config.FromDotfilesDir(dotfilesDir)
	implicitDotIgnore := set.NewFromSlice(config.ImplicitDotIgnore)
	includeFiles := glob_collection.NewGlobCollection(config.IncludeFiles)
	excludeFiles := glob_collection.NewGlobCollection(config.ExcludeFiles)

	for _, file := range files {
		dotfilePath, err := ProcessAddedFile(file, config.TargetDir, config.ImplicitDot, implicitDotIgnore, includeFiles, excludeFiles)
		if err != nil {
			log.Warning("Can't add %s: %v", file, err)
			continue
		}
		err = os.MkdirAll(filepath.Dir(dotfilePath.Str()), 0755)
		if err != nil {
			log.Error("Error creating directory %s: %v", filepath.Dir(dotfilePath.Str()), err)
			continue
		}
		// Hardlink instead of copy, the original file will be replaced on install anyway
		err = os.Link(file, dotfilePath.Str())
		if err != nil {
			log.Error("Error moving %s to %s: %v", file, dotfilePath, err)
		}
	}

	install.Install()
}

func ProcessAddedFile(
	input string,
	targetDir string,
	implicitDot bool,
	implicitDotIgnore set.Set[string],
	includeFiles glob_collection.GlobCollection,
	excludeFiles glob_collection.GlobCollection,
) (RelativePath, error) {
	fileInfo, err := os.Stat(input)
	if err != nil {
		return "", err
	}
	if fileInfo.IsDir() {
		return "", fmt.Errorf("it's a directory. Consider adding %s/**/* instead", input)
	}
	absFile, err := filepath.Abs(input)
	if err != nil {
		return "", fmt.Errorf("error getting absolute path: %v", err)
	}
	if !strings.HasPrefix(absFile, targetDir) {
		return "", fmt.Errorf("it's not inside target directory %s", targetDir)
	}

	targetDirLen := len(targetDir) + len(string(filepath.Separator))
	relPath := NewAbsolutePath(absFile).ExtractRelativePath(targetDirLen)

	if implicitDot && !implicitDotIgnore.Contains(relPath.TopLevelDir()) {
		if relPath.IsHidden() {
			relPath = relPath.Unhide()
		} else {
			return "", fmt.Errorf("its current filename is impossible because implicit_dot is true. Add '%s' to implicit_dot_ignore to fix this", relPath.TopLevelDir())
		}
	}

	if err := checkIsIncluded(relPath, includeFiles, excludeFiles); err != nil {
		return "", err
	}

	return relPath, nil
}

func checkIsIncluded(relPath RelativePath, includeFiles glob_collection.GlobCollection, excludeFiles glob_collection.GlobCollection) error {
	if relPath.Str() == "." {
		return nil
	}
	if excludeFiles.Matches(relPath) && !includeFiles.Matches(relPath) {
		return fmt.Errorf("%s matches exclude_files but is not included in include_files", relPath)
	}
	return checkIsIncluded(relPath.Parent(), includeFiles, excludeFiles)
}
