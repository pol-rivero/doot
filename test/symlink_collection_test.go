package test

import (
	"testing"

	. "github.com/pol-rivero/doot/lib/types"
	"github.com/stretchr/testify/assert"
)

func TestSymlinkCollection(t *testing.T) {
	symlinkCollection := NewSymlinkCollection(42)
	assert.Zero(t, symlinkCollection.Len())
	assert.Empty(t, symlinkCollection.Iter())
	assert.Equal(t, "", symlinkCollection.PrintList())
	assert.Equal(t, "{}", symlinkCollection.ToJson())
	assert.True(t, symlinkCollection.Get(AbsolutePath("/some/path")).IsEmpty())

	symlinkCollection.Add(AbsolutePath("/some/path"), AbsolutePath("/some/target"))
	symlinkCollection.Add(AbsolutePath("/another/path"), AbsolutePath("/another/target"))
	assert.Equal(t, 2, symlinkCollection.Len())
	assert.Equal(t, map[AbsolutePath]AbsolutePath{
		AbsolutePath("/some/path"):    AbsolutePath("/some/target"),
		AbsolutePath("/another/path"): AbsolutePath("/another/target"),
	}, symlinkCollection.Iter())
	assert.Equal(t, AbsolutePath("/some/target"), symlinkCollection.Get(AbsolutePath("/some/path")).Value())

	assert.Equal(t, "/another/path -> /another/target\n/some/path -> /some/target\n", symlinkCollection.PrintList())
	assert.Equal(t, `{"/another/path":"/another/target","/some/path":"/some/target"}`, symlinkCollection.ToJson())
}
