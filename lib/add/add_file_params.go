package add

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/glob_collection"
	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/pol-rivero/doot/lib/utils"
	"github.com/pol-rivero/doot/lib/utils/set"
)

type ProcessAddedFileParams struct {
	crypt             bool
	hostSpecificDir   string
	dotfilesDir       string
	targetDir         string
	implicitDot       bool
	implicitDotIgnore set.Set[string]
	includeFiles      glob_collection.GlobCollection
	excludeFiles      glob_collection.GlobCollection
}

func ProcessAddedFile(input string, params ProcessAddedFileParams) (RelativePath, error) {
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
	cleanAbsFile, err := filepath.Abs(input)
	if err != nil {
		return "", fmt.Errorf("error getting absolute path: %v", err)
	}
	if !strings.HasPrefix(cleanAbsFile, params.targetDir) {
		return "", fmt.Errorf("it's not inside target directory %s", params.targetDir)
	}

	relPath, err := constructRelativePath(cleanAbsFile, params)
	if err != nil {
		return "", fmt.Errorf("error getting relative path: %v", err)
	}

	if params.implicitDot && !params.implicitDotIgnore.Contains(relPath.TopLevelDir()) {
		if relPath.IsHidden() {
			relPath = relPath.Unhide()
		} else {
			return "", fmt.Errorf("its current filename is impossible because implicit_dot is true. Add '%s' to implicit_dot_ignore to fix this", relPath.TopLevelDir())
		}
	}

	if params.crypt && !strings.Contains(relPath.Str(), common.DOOT_CRYPT_EXT) {
		relPath = addDootCryptExtension(relPath)
	}

	if err := checkIsIncluded(relPath, params.includeFiles, params.excludeFiles); err != nil {
		return "", err
	}

	relPath = relPath.AppendLeft(params.hostSpecificDir)

	return relPath, nil
}

func constructRelativePath(absPath string, params ProcessAddedFileParams) (RelativePath, error) {
	relPathStr, err := filepath.Rel(params.targetDir, absPath)
	if err != nil {
		return "", err
	}
	parts := strings.Split(relPathStr, string(filepath.Separator))
	if len(parts) == 0 {
		return "", fmt.Errorf("empty relative path")
	}

	for i := range len(parts) - 1 {
		currentAbsDir := filepath.Join(appendTo(params.dotfilesDir, parts[:i+1])...)
		if stat, err := os.Stat(currentAbsDir); err == nil && stat.IsDir() {
			continue
		}

		log.Info("%s does not exist, checking for crypt directory", currentAbsDir)
		parentDir := filepath.Join(appendTo(params.dotfilesDir, parts[:i])...)
		cryptDir := filepath.Join(parentDir, parts[i]+common.DOOT_CRYPT_EXT)
		if stat, err := os.Stat(cryptDir); err == nil && stat.IsDir() {
			canUseCrypt := params.crypt || utils.RequestInput("Yn", "Do you want to add '%s' inside the existing directory '%s'? It will become encrypted even though you didn't use the --crypt flag. Press N to create a new directory '%s'", absPath, cryptDir, currentAbsDir) == 'y'
			if !canUseCrypt {
				continue
			}
			log.Info("Using crypt directory %s", cryptDir)
			parts[i] = parts[i] + common.DOOT_CRYPT_EXT
		}
	}
	return RelativePath(filepath.Join(parts...)), nil
}

func appendTo(base string, elements []string) []string {
	return append([]string{base}, elements...)
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
	if len(parts) == 0 {
		panic("impossible according to strings.Split spec, since sep is not empty")
	}
	if len(parts) == 1 || (len(parts) == 2 && parts[0] == "") {
		// file.DOOT-CRYPT or .file.DOOT-CRYPT
		parts = append(parts, common.DOOT_CRYPT_EXT_WITHOUT_DOT)
	} else {
		// some.file.DOOT-CRYPT.ext
		parts = append(parts[:len(parts)-1], common.DOOT_CRYPT_EXT_WITHOUT_DOT, parts[len(parts)-1])
	}
	return RelativePath(filepath.Join(dir.Str(), strings.Join(parts, ".")))
}
