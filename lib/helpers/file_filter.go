package helpers

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pol-rivero/doot/lib/config"
	"github.com/pol-rivero/doot/lib/constants"
	"github.com/pol-rivero/doot/lib/log"
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

func ScanDirectory(absolutePath string, filter FileFilter) []string {
	const SEPARATOR_LEN = len(string(filepath.Separator))
	prefixLen := len(absolutePath) + SEPARATOR_LEN
	files := make([]string, 0, 64)
	scanDirectoryRecursive(filter, &files, prefixLen, absolutePath)
	return files
}

func scanDirectoryRecursive(filter FileFilter, files *[]string, prefixLen int, absolutePath string) {
	entries, err := os.ReadDir(absolutePath)
	if err != nil {
		log.Error("Error reading directory %s: %v", absolutePath, err)
		return
	}
	for _, entry := range entries {
		entryName := entry.Name()
		entryAbsPath := filepath.Join(absolutePath, entryName)
		entryRelativePath := entryAbsPath[prefixLen:]
		if filter.isExcluded(entryRelativePath, entryName) {
			continue
		}

		if entry.IsDir() {
			scanDirectoryRecursive(filter, files, prefixLen, entryAbsPath)
		} else {
			*files = append(*files, entryRelativePath)
		}
	}
}

func (f *FileFilter) isExcluded(path string, fileName string) bool {
	return f.matchesExcludePattern(path, fileName) && !f.IncludeGlobs.Matches(path)
}

func (f *FileFilter) matchesExcludePattern(path string, fileName string) bool {
	if f.IgnoreHidden && fileName[0] == '.' {
		return true
	}
	if f.IgnoreDootCrypt && strings.Contains(fileName, constants.DOOT_CRYPT_EXT) {
		return true
	}
	return f.ExcludeGlobs.Matches(path)
}
