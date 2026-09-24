package common

import (
	"path/filepath"
	"strings"
)

const separator = string(filepath.Separator)

// Must agree with the git attributes written by 'doot crypt init'. A file is encrypted if its
// name, or the name of any of its parent directories, matches *.doot-crypt or *.doot-crypt.*
func IsCryptName(name string) bool {
	return strings.HasSuffix(name, DOOT_CRYPT_EXT) || strings.Contains(name, DOOT_CRYPT_EXT+".")
}

func IsCryptPath(path string) bool {
	for part := range strings.SplitSeq(path, separator) {
		if IsCryptName(part) {
			return true
		}
	}
	return false
}

// Removes only the occurrences of .doot-crypt that make the path encrypted (the ones at the end
// of a name or followed by a dot). Other occurrences, like in "token.doot-crypt_old", are kept.
func RemoveCryptExt(path string) string {
	if !strings.Contains(path, DOOT_CRYPT_EXT) {
		return path
	}
	parts := strings.Split(path, separator)
	for i, part := range parts {
		parts[i] = removeCryptExtFromName(part)
	}
	return strings.Join(parts, separator)
}

func removeCryptExtFromName(name string) string {
	var result strings.Builder
	for {
		index := strings.Index(name, DOOT_CRYPT_EXT)
		if index == -1 {
			result.WriteString(name)
			return result.String()
		}
		rest := name[index+len(DOOT_CRYPT_EXT):]
		if rest == "" || rest[0] == '.' {
			result.WriteString(name[:index])
		} else {
			result.WriteString(name[:index+len(DOOT_CRYPT_EXT)])
		}
		name = rest
	}
}
