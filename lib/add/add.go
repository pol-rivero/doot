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
	"github.com/pol-rivero/doot/lib/crypt"
	"github.com/pol-rivero/doot/lib/install"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/pol-rivero/doot/lib/utils/set"
)

func Add(files []string, isCrypt bool, isHostSpecific bool) {
	dotfilesDir := common.FindDotfilesDir()
	config := config.FromDotfilesDir(dotfilesDir)
	params := ProcessAddedFileParams{
		crypt:             isCrypt,
		hostSpecific:      isHostSpecific,
		targetDir:         config.TargetDir,
		implicitDot:       config.ImplicitDot,
		implicitDotIgnore: set.NewFromSlice(config.ImplicitDotIgnore),
		includeFiles:      glob_collection.NewGlobCollection(config.IncludeFiles),
		excludeFiles:      glob_collection.NewGlobCollection(config.ExcludeFiles),
	}

	if isCrypt && !crypt.GitCryptIsInitialized(dotfilesDir) {
		log.Error("Can't add private files with --crypt flag because repository is not initialized. Run 'doot crypt init' first.")
		return
	}

	for _, file := range files {
		dotfileRelativePath, err := processAddedFile(file, params)
		if err != nil {
			log.Error("Can't add %s: %v", file, err)
			continue
		}
		dotfilePath := dotfilesDir.JoinPath(dotfileRelativePath)
		err = os.MkdirAll(filepath.Dir(dotfilePath.Str()), 0755)
		if err != nil {
			log.Error("Error creating directory %s: %v", filepath.Dir(dotfilePath.Str()), err)
			continue
		}
		// Hardlink instead of copy, the original file will be replaced on install anyway
		err = os.Link(file, dotfilePath.Str())
		if err != nil {
			log.Error("Error moving %s to %s: %v", file, dotfilePath, err)
		} else {
			log.Info("Created hardlink %s -> %s", file, dotfilePath)
		}
	}

	log.Info("Files have been copied to the dotfiles directory, now running 'install'...")
	install.Install()
}

type ProcessAddedFileParams struct {
	crypt             bool
	hostSpecific      bool
	targetDir         string
	implicitDot       bool
	implicitDotIgnore set.Set[string]
	includeFiles      glob_collection.GlobCollection
	excludeFiles      glob_collection.GlobCollection
}

func processAddedFile(input string, params ProcessAddedFileParams) (RelativePath, error) {
	fileInfo, err := os.Stat(input)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("this file does not exist")
		}
		return "", err
	}
	if fileInfo.IsDir() {
		return "", fmt.Errorf("it's a directory. Consider adding %s/**/* instead", input)
	}
	absFile, err := filepath.Abs(input)
	if err != nil {
		return "", fmt.Errorf("error getting absolute path: %v", err)
	}
	if !strings.HasPrefix(absFile, params.targetDir) {
		return "", fmt.Errorf("it's not inside target directory %s", params.targetDir)
	}

	targetDirLen := len(params.targetDir) + len(string(filepath.Separator))
	relPath := NewAbsolutePath(absFile).ExtractRelativePath(targetDirLen)

	if params.implicitDot && !params.implicitDotIgnore.Contains(relPath.TopLevelDir()) {
		if relPath.IsHidden() {
			relPath = relPath.Unhide()
		} else {
			return "", fmt.Errorf("its current filename is impossible because implicit_dot is true. Add '%s' to implicit_dot_ignore to fix this", relPath.TopLevelDir())
		}
	}

	if err := checkIsIncluded(relPath, params.includeFiles, params.excludeFiles); err != nil {
		return "", err
	}

	if params.crypt {
		relPath = addDootCryptExtension(relPath)
	}

	if params.hostSpecific {
		// TODO
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

func addDootCryptExtension(relPath RelativePath) RelativePath {
	dir, file := relPath.Split()
	parts := strings.Split(file, ".")
	if len(parts) > 1 {
		parts = append(parts[:len(parts)-1], common.DOOT_CRYPT_EXT_WITHOUT_DOT, parts[len(parts)-1])
	} else {
		parts = append(parts, common.DOOT_CRYPT_EXT_WITHOUT_DOT)
	}
	return RelativePath(filepath.Join(dir.Str(), strings.Join(parts, ".")))
}
