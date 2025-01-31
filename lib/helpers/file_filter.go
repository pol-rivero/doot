package helpers

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
	"github.com/pol-rivero/doot/lib/config"
)

type FileFilter struct {
	IgnoreHidden    bool
	IgnoreDootCrypt bool
	ExcludeGlobs    []glob.Glob
	IncludeGlobs    []glob.Glob
}

func CreateFilter(config *config.Config, ignoreDootCrypt bool) FileFilter {
	ignoreHidden := false
	var excludeGlobs = make([]glob.Glob, 0, len(config.ExcludeFiles))
	for _, excludePattern := range config.ExcludeFiles {
		if excludePattern == "**/.*" {
			ignoreHidden = true
			continue
		}
		g, err := glob.Compile(excludePattern, filepath.Separator)
		if err != nil {
			log.Fatalf("Ignoring invalid exclude pattern '%s': %v", excludePattern, err)
			continue
		}
		excludeGlobs = append(excludeGlobs, g)
	}

	var includeGlobs = make([]glob.Glob, 0, len(config.IncludeFiles))
	for _, includePattern := range config.IncludeFiles {
		g, err := glob.Compile(includePattern, filepath.Separator)
		if err != nil {
			log.Fatalf("Ignoring invalid include pattern '%s': %v", includePattern, err)
			continue
		}
		includeGlobs = append(includeGlobs, g)
	}

	return FileFilter{
		IgnoreHidden:    ignoreHidden,
		IgnoreDootCrypt: ignoreDootCrypt,
		ExcludeGlobs:    excludeGlobs,
		IncludeGlobs:    includeGlobs,
	}
}

func ScanDirectory(absolutePath string, filter FileFilter) []string {
	const SEPARATOR_LEN = len(string(filepath.Separator))
	prefixLen := len(absolutePath) + SEPARATOR_LEN
	files := make([]string, 0, 32)
	scanDirectoryRecursive(filter, &files, prefixLen, absolutePath)
	return files
}

func scanDirectoryRecursive(filter FileFilter, files *[]string, prefixLen int, absolutePath string) {
	entries, err := os.ReadDir(absolutePath)
	if err != nil {
		log.Fatalf("Error reading directory %s: %v", absolutePath, err)
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
	return matchesExcludePattern(f, path, fileName) && !matchesIncludePattern(f, path)
}

func matchesExcludePattern(filter *FileFilter, path string, fileName string) bool {
	if filter.IgnoreHidden && fileName[0] == '.' {
		return true
	}
	if filter.IgnoreDootCrypt && strings.Contains(fileName, DOOT_CRYPT_EXT) {
		return true
	}
	for _, excludeGlob := range filter.ExcludeGlobs {
		if excludeGlob.Match(path) {
			return true
		}
	}
	return false
}

func matchesIncludePattern(filter *FileFilter, path string) bool {
	for _, includeGlob := range filter.IncludeGlobs {
		if includeGlob.Match(path) {
			return true
		}
	}
	return false
}
