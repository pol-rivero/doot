package install

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pol-rivero/doot/lib/config"
	"github.com/pol-rivero/doot/lib/constants"
	"github.com/pol-rivero/doot/lib/log"
	. "github.com/pol-rivero/doot/lib/types"
)

type FileFilter struct {
	IgnoreHidden    bool
	IgnoreDootCrypt bool
	ExcludeGlobs    GlobCollection
	IncludeGlobs    GlobCollection
}

func CreateFilter(config *config.Config, ignoreDootCrypt bool) FileFilter {
	ignoreHidden := false
	newExcludeFiles := make([]string, 0, len(config.ExcludeFiles))
	for _, excludePattern := range config.ExcludeFiles {
		if excludePattern == "**/.*" {
			ignoreHidden = true
		} else {
			newExcludeFiles = append(newExcludeFiles, excludePattern)
		}
	}
	return FileFilter{
		IgnoreHidden:    ignoreHidden,
		IgnoreDootCrypt: ignoreDootCrypt,
		ExcludeGlobs:    NewGlobCollection(newExcludeFiles),
		IncludeGlobs:    NewGlobCollection(config.IncludeFiles),
	}
}

func ScanDirectory(dir AbsolutePath, filter FileFilter) []RelativePath {
	const SEPARATOR_LEN = len(string(filepath.Separator))
	prefixLen := len(dir) + SEPARATOR_LEN
	files := make([]RelativePath, 0, 64)
	scanDirectoryRecursive(filter, &files, prefixLen, dir)
	return files
}

func scanDirectoryRecursive(filter FileFilter, result *[]RelativePath, prefixLen int, scanPath AbsolutePath) {
	entries, err := os.ReadDir(scanPath.Str())
	if err != nil {
		log.Error("Error reading directory %s: %v", scanPath, err)
		return
	}
	for _, entry := range entries {
		entryName := entry.Name()
		entryPath := scanPath.Join(entryName)
		entryRelativePath := entryPath.ExtractRelativePath(prefixLen)
		if filter.isExcluded(entryRelativePath, entryName) {
			continue
		}

		if entry.IsDir() {
			scanDirectoryRecursive(filter, result, prefixLen, entryPath)
		} else {
			*result = append(*result, entryRelativePath)
		}
	}
}

func (f *FileFilter) isExcluded(path RelativePath, fileName string) bool {
	return f.matchesExcludePattern(path, fileName) && !f.IncludeGlobs.Matches(path)
}

func (f *FileFilter) matchesExcludePattern(path RelativePath, fileName string) bool {
	if f.IgnoreHidden && fileName[0] == '.' {
		return true
	}
	if f.IgnoreDootCrypt && strings.Contains(fileName, constants.DOOT_CRYPT_EXT) {
		return true
	}
	return f.ExcludeGlobs.Matches(path)
}
