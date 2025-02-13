package types

import (
	"path/filepath"
	"strings"
)

type RelativePath string
type AbsolutePath string

func (rp RelativePath) Str() string {
	return string(rp)
}

func (rp RelativePath) MakeAbsolute(baseDir AbsolutePath) AbsolutePath {
	return baseDir.Join(rp.Str())
}

func (rp RelativePath) Replace(substring, replacement string) RelativePath {
	return RelativePath(strings.ReplaceAll(rp.Str(), substring, replacement))
}

func (rp RelativePath) RemoveBaseDir(baseDirLen int) RelativePath {
	return RelativePath(rp.Str()[baseDirLen:])
}

func (rp RelativePath) TopLevelDir() string {
	firstSeparatorIndex := strings.IndexRune(rp.Str(), filepath.Separator)
	if firstSeparatorIndex == -1 {
		return rp.Str()
	}
	return rp.Str()[:firstSeparatorIndex]
}

func (rp RelativePath) Split() (RelativePath, string) {
	dir, file := filepath.Split(rp.Str())
	return RelativePath(dir), file
}

func (rp RelativePath) IsHidden() bool {
	return strings.HasPrefix(rp.Str(), ".")
}

func (rp RelativePath) Unhide() RelativePath {
	if !rp.IsHidden() {
		return rp
	}
	return RelativePath(rp.Str()[1:])
}

func (ap RelativePath) Parent() RelativePath {
	return RelativePath(filepath.Dir(ap.Str()))
}

func (ap AbsolutePath) Str() string {
	return string(ap)
}

func NewAbsolutePath(path string) AbsolutePath {
	if !filepath.IsAbs(path) {
		panic("Attempted to create AbsolutePath from non-absolute path")
	}
	return AbsolutePath(path)
}

func (ap AbsolutePath) Join(other string) AbsolutePath {
	return AbsolutePath(filepath.Join(ap.Str(), other))
}

func (ap AbsolutePath) JoinPath(other RelativePath) AbsolutePath {
	return ap.Join(other.Str())
}

func (ap AbsolutePath) ExtractRelativePath(baseDirLen int) RelativePath {
	return RelativePath(ap.Str()[baseDirLen:])
}

func (ap AbsolutePath) Parent() AbsolutePath {
	return AbsolutePath(filepath.Dir(ap.Str()))
}
