package types

import "github.com/pol-rivero/doot/lib/utils/optional"

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

func (sc *SymlinkCollection) Len() int {
	return len(sc.links)
}

func (sc *SymlinkCollection) Iter() map[AbsolutePath]AbsolutePath {
	return sc.links
}
