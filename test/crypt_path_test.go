package test

import (
	"testing"

	"github.com/pol-rivero/doot/lib/common"
	"github.com/stretchr/testify/assert"
)

// Must agree with the git attributes written by 'doot crypt init':
// *.doot-crypt, *.doot-crypt.*, *.doot-crypt/** and *.doot-crypt.*/**
func TestCryptPath_IsCryptName(t *testing.T) {
	cryptNames := []string{
		"file.doot-crypt",
		"file.doot-crypt.txt",
		".file.doot-crypt",
		".file.doot-crypt.txt",
		"some.file.doot-crypt.ext",
		"dir.doot-crypt.d",
		"a.doot-crypt_x.doot-crypt",
	}
	for _, name := range cryptNames {
		assert.True(t, common.IsCryptName(name), "Expected %s to be a crypt name", name)
	}

	notCryptNames := []string{
		"file",
		"file.txt",
		"token.doot-crypt_old",
		"keys.doot-crypted",
		"file.doot-crypt-backup",
		"file-doot-crypt",
		"file.doot-cryptx.txt",
	}
	for _, name := range notCryptNames {
		assert.False(t, common.IsCryptName(name), "Expected %s not to be a crypt name", name)
	}
}

func TestCryptPath_IsCryptPath(t *testing.T) {
	assert.True(t, common.IsCryptPath("file.doot-crypt"))
	assert.True(t, common.IsCryptPath("dir/file.doot-crypt.txt"))
	assert.True(t, common.IsCryptPath("dir.doot-crypt/file"))
	assert.True(t, common.IsCryptPath("dir.doot-crypt.d/nested/file"))

	assert.False(t, common.IsCryptPath("file"))
	assert.False(t, common.IsCryptPath("notes.doot-crypted"))
	assert.False(t, common.IsCryptPath("old.doot-crypt_bak/key"))
	assert.False(t, common.IsCryptPath("dir/token.doot-crypt_old"))
}

func TestCryptPath_RemoveCryptExt(t *testing.T) {
	expected := map[string]string{
		"file":                                "file",
		"file.doot-crypt":                     "file",
		"file.doot-crypt.txt":                 "file.txt",
		".file.doot-crypt":                    ".file",
		"dirA.doot-crypt/dirB.doot-crypt.d/f": "dirA/dirB.d/f",
		"a.doot-crypt.doot-crypt.b":           "a.b",
		"token.doot-crypt_old":                "token.doot-crypt_old",
		"keys.doot-crypted/id":                "keys.doot-crypted/id",
		"notes.doot-crypt.doot-crypted":       "notes.doot-crypted",
		"a.doot-crypt_x.doot-crypt":           "a.doot-crypt_x",
	}
	for input, output := range expected {
		assert.Equal(t, output, common.RemoveCryptExt(input), "Unexpected result for %s", input)
	}
}
