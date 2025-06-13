package types

import (
	"encoding/json"
	"slices"
	"strings"

	"github.com/pol-rivero/doot/lib/utils/optional"
)

type SymlinkCollection struct {
	// link path -> link content (target)
	links map[AbsolutePath]AbsolutePath
}

func NewSymlinkCollection(capacity int) SymlinkCollection {
	return SymlinkCollection{make(map[AbsolutePath]AbsolutePath, capacity)}
}

func (sc *SymlinkCollection) Add(linkPath, linkContent AbsolutePath) {
	sc.links[linkPath] = linkContent
}

func (sc *SymlinkCollection) Get(linkPath AbsolutePath) optional.Optional[AbsolutePath] {
	linkContent, exists := sc.links[linkPath]
	if exists {
		return optional.Of(linkContent)
	}
	return optional.Empty[AbsolutePath]()
}

func (sc *SymlinkCollection) Remove(linkPath AbsolutePath) {
	delete(sc.links, linkPath)
}

func (sc *SymlinkCollection) Len() int {
	return len(sc.links)
}

func (sc *SymlinkCollection) Iter() map[AbsolutePath]AbsolutePath {
	return sc.links
}

func (sc *SymlinkCollection) PrintList() string {
	paths := make([]string, 0, len(sc.links))
	for path := range sc.links {
		paths = append(paths, path.Str())
	}
	slices.Sort(paths)

	var sb strings.Builder
	for _, pathStr := range paths {
		path := AbsolutePath(pathStr)
		content := sc.links[path]
		sb.WriteString(path.Str())
		sb.WriteString(" -> ")
		sb.WriteString(content.Str())
		sb.WriteString("\n")
	}
	return sb.String()
}

func (sc *SymlinkCollection) ToJson() string {
	jsonBytes, err := json.Marshal(sc.links)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}
