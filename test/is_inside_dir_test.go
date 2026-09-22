package test

import (
	"testing"

	"github.com/pol-rivero/doot/lib/common"
	"github.com/stretchr/testify/assert"
)

func TestIsInsideDir(t *testing.T) {
	assert.True(t, common.IsInsideDir("/a/dots", "/a/dots/file"))
	assert.True(t, common.IsInsideDir("/a/dots", "/a/dots/dir/file"))
	assert.True(t, common.IsInsideDir("/a/dots/", "/a/dots/file"))
	assert.True(t, common.IsInsideDir("/a/dots", "/a/dots/..file"))

	assert.False(t, common.IsInsideDir("/a/dots", "/a/dots"))
	assert.False(t, common.IsInsideDir("/a/dots", "/a/dots/"))
	assert.False(t, common.IsInsideDir("/a/dots", "/a/dots-old/file"))
	assert.False(t, common.IsInsideDir("/a/dots", "/a/dotsfile"))
	assert.False(t, common.IsInsideDir("/a/dots", "/a/dots/../dots-old/file"))
	assert.False(t, common.IsInsideDir("/a/dots", "/a"))
	assert.False(t, common.IsInsideDir("/a/dots", "dots/file"))
}
