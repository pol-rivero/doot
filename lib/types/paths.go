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
